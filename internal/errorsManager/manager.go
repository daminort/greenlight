package errorsManager

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"greenlight.damian.net/internal/pkg/envelopes"
	"greenlight.damian.net/internal/pkg/payloads"
	"greenlight.damian.net/internal/pkg/validator"
)

type ErrorsManager struct {
	Logger *slog.Logger
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func New(logger *slog.Logger) *ErrorsManager {
	return &ErrorsManager{
		Logger: logger,
	}
}

func (m *ErrorsManager) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	m.Logger.Error(err.Error(),
		"method", method,
		"uri", uri,
	)
}

func (m *ErrorsManager) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envelope := envelopes.New("error", message)

	err := payloads.WriteJSON(w, status, envelope, nil)
	if err != nil {
		m.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (m *ErrorsManager) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	m.LogError(r, err)

	message := "An internal server error has occurred."
	m.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (m *ErrorsManager) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found."
	m.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (m *ErrorsManager) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The requested method %s is not allowed.", r.Method)
	m.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (m *ErrorsManager) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	m.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (m *ErrorsManager) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors validator.ValidationErrors) {
	m.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (m *ErrorsManager) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprint("unable to update the record due to an edit conflict, please try again")
	m.ErrorResponse(w, r, http.StatusConflict, message)
}
