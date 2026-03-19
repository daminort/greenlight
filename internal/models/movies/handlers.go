package movies

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/pkg/envelopes"
	"greenlight.damian.net/internal/pkg/filters"
	"greenlight.damian.net/internal/pkg/payloads"
	"greenlight.damian.net/internal/pkg/queries"
	"greenlight.damian.net/internal/pkg/requests"
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

func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {
	fiParams := filters.InitParams{
		SearchKey:   "title",
		SortDefault: "title",
		Columns: []string{
			"title", "-title",
			"year", "-year",
			"runtime", "-runtime",
		},
	}

	fils := filters.New(r.URL.Query(), fiParams)
	fErrors := fils.Validate()
	if len(fErrors) != 0 {
		envelop := envelopes.New("errors", fErrors)
		err := payloads.WriteJSON(w, http.StatusBadRequest, envelop, nil)
		if err != nil {
			h.ErrorManager.ServerErrorResponse(w, r, err)
		}

		return
	}

	query := queries.New(r.URL.Query())
	params := GetMoviesParams{
		Genres:  query.ReadStrings("genres", []string{}),
		Filters: fils,
	}

	movies, meta, err := h.Service.GetList(params)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"movies": movies,
		"meta":   meta,
	}

	envelop := envelopes.NewPack(data)
	err = payloads.WriteJSON(w, http.StatusOK, envelop, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ReadParamInt(r, "id")
	if err != nil {
		h.ErrorManager.NotFoundResponse(w, r)
		return
	}

	movie, err := h.Service.Get(id)
	if err != nil {
		if errors.Is(err, errorsManager.ErrRecordNotFound) {
			h.ErrorManager.NotFoundResponse(w, r)
			return
		}

		h.ErrorManager.ServerErrorResponse(w, r, err)
		return
	}

	envelope := envelopes.New("movie", movie)

	err = payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateMoviePayload

	err := payloads.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorManager.BadRequestResponse(w, r, err)
		return
	}

	movie := &Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	form := ValidateMovie(movie)
	if !form.IsValid() {
		h.ErrorManager.FailedValidationResponse(w, r, form.Errors)
		return
	}

	err = h.Service.Create(movie)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	envelop := envelopes.New("movie", movie)
	err = payloads.WriteJSON(w, http.StatusCreated, envelop, headers)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var input UpdateMoviePayload

	id, err := requests.ReadParamInt(r, "id")
	if err != nil {
		h.ErrorManager.NotFoundResponse(w, r)
		return
	}

	movie, err := h.Service.Get(id)
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

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	form := ValidateMovie(movie)
	if !form.IsValid() {
		h.ErrorManager.FailedValidationResponse(w, r, form.Errors)
		return
	}

	err = h.Service.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, errorsManager.ErrEditConflict):
			h.ErrorManager.EditConflictResponse(w, r)
		default:
			h.ErrorManager.ServerErrorResponse(w, r, err)
		}
		return
	}

	envelope := envelopes.New("movie", movie)
	err = payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ReadParamInt(r, "id")
	if err != nil {
		h.ErrorManager.NotFoundResponse(w, r)
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		if errors.Is(err, errorsManager.ErrRecordNotFound) {
			h.ErrorManager.NotFoundResponse(w, r)
			return
		}
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}

	envelope := envelopes.New("response", map[string]string{"status": "ok"})
	err = payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}
