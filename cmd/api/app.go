package main

import (
	"log/slog"

	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/middlewares"
	"greenlight.damian.net/internal/models/health"
	"greenlight.damian.net/internal/models/movies"
)

type Application struct {
	Config       *config.Config
	Logger       *slog.Logger
	ErrorManager *errorsManager.ErrorsManager
	Middlewares  *middlewares.Middlewares
	Movies       *movies.Handlers
	Health       *health.Handlers
}
