package users

import (
	"errors"
	"net/http"

	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/pkg/envelopes"
	"greenlight.damian.net/internal/pkg/payloads"
	"greenlight.damian.net/internal/pkg/queries"
	"greenlight.damian.net/internal/pkg/requests"
	"greenlight.damian.net/internal/pkg/validator"
)

type Handlers struct {
	Service      ServiceInstance
	ErrorManager *errorsManager.ErrorsManager
}

func NewHandlers(s ServiceInstance, em *errorsManager.ErrorsManager) *Handlers {
	return &Handlers{
		Service:      s,
		ErrorManager: em,
	}
}

func (h *Handlers) GetByEmail(w http.ResponseWriter, r *http.Request) {
	query := queries.New(r.URL.Query())
	email := query.ReadString("email", "")

	v := validator.New()
	v.Check(validator.IsEmail(email), "email", "must be a valid email address")
	if !v.IsValid() {
		h.ErrorManager.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.Service.GetByEmail(email)
	if err != nil {
		if errors.Is(err, errorsManager.ErrRecordNotFound) {
			h.ErrorManager.NotFoundResponse(w, r)
			return
		}

		h.ErrorManager.ServerErrorResponse(w, r, err)
		return
	}

	envelope := envelopes.New("user", user)

	err = payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateUserPayload

	err := payloads.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorManager.BadRequestResponse(w, r, err)
		return
	}

	pwd := Password{}
	err = pwd.Set(input.Pwd)
	if err != nil {
		h.ErrorManager.BadRequestResponse(w, r, err)
		return
	}

	user := &User{
		Name:      input.Name,
		Email:     input.Email,
		Pwd:       pwd,
		Activated: false,
	}

	form := ValidateUser(user)
	if !form.IsValid() {
		h.ErrorManager.FailedValidationResponse(w, r, form.Errors)
		return
	}

	err = h.Service.Create(user)
	if err != nil {
		if errors.Is(err, errorsManager.ErrDuplicateEmail) {
			h.ErrorManager.EmailDuplicatedResponse(w, r)
			return
		}

		h.ErrorManager.ServerErrorResponse(w, r, err)
		return
	}

	envelop := envelopes.New("user", user)
	err = payloads.WriteJSON(w, http.StatusCreated, envelop, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var input UpdateUserPayload

	id, err := requests.ReadParamInt(r, "id")
	if err != nil {
		h.ErrorManager.NotFoundResponse(w, r)
		return
	}

	user, err := h.Service.Get(id)
	if err != nil {
		if errors.Is(err, errorsManager.ErrRecordNotFound) {
			h.ErrorManager.NotFoundResponse(w, r)
			return
		}
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}

	err = payloads.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorManager.BadRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Activated != nil {
		user.Activated = *input.Activated
	}

	if input.Pwd != nil {
		pwd := Password{}
		err = pwd.Set(*input.Pwd)
		if err != nil {
			h.ErrorManager.BadRequestResponse(w, r, err)
			return
		}
		user.Pwd = pwd
	}

	form := ValidateUser(user)
	if !form.IsValid() {
		h.ErrorManager.FailedValidationResponse(w, r, form.Errors)
		return
	}

	err = h.Service.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, errorsManager.ErrEditConflict):
			h.ErrorManager.EditConflictResponse(w, r)
		case errors.Is(err, errorsManager.ErrDuplicateEmail):
			h.ErrorManager.EmailDuplicatedResponse(w, r)
		default:
			h.ErrorManager.ServerErrorResponse(w, r, err)
		}
		return
	}

	envelope := envelopes.New("user", user)
	err = payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}
