package main

import (
	"fmt"
	"net/http"

	"greenlight.damian.net/internal/utils"
)

func (app *Application) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(),
		"method", method,
		"uri", uri,
	)
}

func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envelope := utils.NewEnvelope("error", message)

	err := utils.WriteJSON(w, status, envelope, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	message := "An internal server error has occurred."
	app.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found."
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The requested method %s is not allowed.", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}
