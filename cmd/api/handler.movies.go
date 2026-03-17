package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"greenlight.damian.net/internal/filters"
	"greenlight.damian.net/internal/models"
	"greenlight.damian.net/internal/queries"
	"greenlight.damian.net/internal/utils"
	"greenlight.damian.net/internal/validator"
)

// Handlers

func (app *Application) getMovies(w http.ResponseWriter, r *http.Request) {
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
		envelop := utils.NewEnvelope("errors", fErrors)
		err := utils.WriteJSON(w, http.StatusBadRequest, envelop, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}

		return
	}

	query := queries.New(r.URL.Query())
	params := models.GetMoviesParams{
		Genres:  query.ReadStrings("genres", []string{}),
		Filters: fils,
	}

	movies, meta, err := app.Models.Movies.GetMovies(params)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"movies": movies,
		"meta":   meta,
	}

	envelop := utils.NewEnvelope("result", data)
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
	var input models.CreateMoviePayload

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	movie := models.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	form := validateMovie(&movie)
	if !form.IsValid() {
		app.FailedValidationResponse(w, r, form.Errors)
		return
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
	var input models.UpdateMoviePayload

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

	form := validateMovie(movie)
	if !form.IsValid() {
		app.FailedValidationResponse(w, r, form.Errors)
		return
	}

	err = app.Models.Movies.UpdateMovie(movie)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.EditConflictResponse(w, r)
		default:
			app.ServerErrorResponse(w, r, err)
		}
		return
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

func validateMovie(m *models.Movie) *validator.Validator {
	v := validator.New()

	v.Check(validator.NotBlank(m.Title), "title", "must be provided")
	v.Check(validator.MaxChars(m.Title, 50), "title", "must not be more than 50 characters")

	v.Check(validator.NotZero(m.Year), "year", "must be provided")
	v.Check(validator.GreaterThan(m.Year, 1887), "year", "must be greater than or equal to 1888")
	v.Check(validator.LessThan(m.Year, int(time.Now().Year())+1), "year", "must not be in the future")

	v.Check(validator.NotZero(int(m.Runtime)), "runtime", "must be provided")
	v.Check(validator.GreaterThan(int(m.Runtime), 0), "runtime", "must be positive")

	v.Check(validator.NotNil(m.Genres), "genres", "must be provided")
	v.Check(validator.GreaterThan(len(m.Genres), 0), "genres", "must contain at least 1 genre")
	v.Check(validator.LessThan(len(m.Genres), 6), "genres", "must not contain more than 5 genres")
	v.Check(validator.IsUnique(m.Genres), "genres", "must not contain duplicate values")

	return v
}
