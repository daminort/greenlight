package main

import (
	"log/slog"

	"greenlight.damian.net/cmd/api/config"
)

type Application struct {
	Config *config.Config
	Logger *slog.Logger
}
