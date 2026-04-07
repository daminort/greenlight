package main

import (
	"log/slog"

	bgManager "greenlight.damian.net/internal/bg_manager"
	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/mailer"
	"greenlight.damian.net/internal/middlewares"
	"greenlight.damian.net/internal/models/health"
	"greenlight.damian.net/internal/models/movies"
	"greenlight.damian.net/internal/models/users"
)

type Application struct {
	Config       *config.Config
	Logger       *slog.Logger
	ErrorManager *errorsManager.ErrorsManager
	BgManager    *bgManager.BgManager
	Mailer       *mailer.Mailer
	Middlewares  *middlewares.Middlewares
	Movies       *movies.Handlers
	Users        *users.Handlers
	Health       *health.Handlers
}
