package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"greenlight.damian.net/internal/models"
	"greenlight.damian.net/internal/utils"
	"greenlight.damian.net/internal/validator"
)

type upsertMoviePayload struct {
	Title   string         `json:"title"`
	Year    int            `json:"year"`
	Runtime models.Runtime `json:"runtime"`
	Genres  []string       `json:"genres"`
}

// Handlers

func (app *Application) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.Models.Movies.GetMovies()
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	envelop := utils.NewEnvelope("movies", movies)
	err = utils.WriteJSON(w, http.StatusOK, envelop, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func (app *Application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamInt(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	movie, err := app.Models.Movies.GetMovie(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.NotFoundResponse(w, r)
			return
		}

		app.ServerErrorResponse(w, r, err)
		return
	}

	envelope := utils.NewEnvelope("movie", movie)

	err = utils.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func (app *Application) createMovie(w http.ResponseWriter, r *http.Request) {
	var input upsertMoviePayload

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	form := validateUpsertMoviePayload(input)
	if !form.IsValid() {
		app.FailedValidationResponse(w, r, form.Errors)
		return
	}

	movie := models.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	err = app.Models.Movies.InsertMovie(&movie)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	envelop := utils.NewEnvelope("movie", movie)
	err = utils.WriteJSON(w, http.StatusCreated, envelop, headers)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func (app *Application) updateMovie(w http.ResponseWriter, r *http.Request) {
	var input upsertMoviePayload

	id, err := utils.ReadParamInt(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	movie, err := app.Models.Movies.GetMovie(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.NotFoundResponse(w, r)
			return
		}
		app.ServerErrorResponse(w, r, err)
	}

	err = utils.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	form := validateUpsertMoviePayload(input)
	if !form.IsValid() {
		app.FailedValidationResponse(w, r, form.Errors)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	err = app.Models.Movies.UpdateMovie(movie)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.NotFoundResponse(w, r)
			return
		}
		app.ServerErrorResponse(w, r, err)
	}

	envelope := utils.NewEnvelope("movie", movie)
	err = utils.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func (app *Application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamInt(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	err = app.Models.Movies.DeleteMovie(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.NotFoundResponse(w, r)
			return
		}
		app.ServerErrorResponse(w, r, err)
	}

	envelope := utils.NewEnvelope("response", map[string]string{"status": "ok"})
	err = utils.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

// Validators

func validateUpsertMoviePayload(p upsertMoviePayload) *validator.Validator {
	v := validator.New()

	v.Check(validator.NotBlank(p.Title), "title", "must be provided")
	v.Check(validator.MaxChars(p.Title, 50), "title", "must not be more than 50 characters")

	v.Check(validator.NotZero(p.Year), "year", "must be provided")
	v.Check(validator.GreaterThan(p.Year, 1887), "year", "must be greater than or equal to 1888")
	v.Check(validator.LessThan(p.Year, int(time.Now().Year())+1), "year", "must not be in the future")

	v.Check(validator.NotZero(int(p.Runtime)), "runtime", "must be provided")
	v.Check(validator.GreaterThan(int(p.Runtime), 0), "runtime", "must be positive")

	v.Check(validator.NotNil(p.Genres), "genres", "must be provided")
	v.Check(validator.GreaterThan(len(p.Genres), 0), "genres", "must contain at least 1 genre")
	v.Check(validator.LessThan(len(p.Genres), 6), "genres", "must not contain more than 5 genres")
	v.Check(validator.IsUnique(p.Genres), "genres", "must not contain duplicate values")

	return v
}
