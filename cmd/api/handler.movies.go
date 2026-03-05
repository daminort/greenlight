package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.damian.net/internal/models"
	"greenlight.damian.net/internal/utils"
)

func (app *Application) createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "createMovie")
}

func (app *Application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamInt(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	movie := models.Movie{
		ID:        id,
		Title:     "Rambo",
		Year:      1985,
		Runtime:   92,
		Genres:    []string{"action", "war"},
		Version:   1,
		CreatedAt: time.Now(),
	}

	envelope := utils.NewEnvelope("movie", movie)

	err = utils.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
