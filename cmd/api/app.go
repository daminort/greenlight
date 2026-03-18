package main

import (
	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/errorsManager"
	"greenlight.damian.net/internal/middlewares"
	"greenlight.damian.net/internal/models/health"
	"greenlight.damian.net/internal/models/movies"
)

type Application struct {
	Config       *config.Config
	ErrorManager *errorsManager.ErrorsManager
	Middlewares  *middlewares.Middlewares
	Movies       *movies.Handlers
	Health       *health.Handlers
}
