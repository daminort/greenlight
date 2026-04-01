package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.ErrorManager.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.ErrorManager.MethodNotAllowedResponse)

	// health
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.Health.Check)

	// movies
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.Movies.GetList)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.Movies.Create)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.Movies.Get)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.Movies.Update)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.Movies.Delete)

	// users
	router.HandlerFunc(http.MethodGet, "/v1/users", app.Users.GetByEmail)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.Users.Create)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", app.Users.Update)

	return app.Middlewares.RecoverPanic(app.Middlewares.RateLimit(router))
}
