package main

import (
	"log/slog"

	"greenlight.damian.net/cmd/api/config"
	"greenlight.damian.net/cmd/api/database"
)

type Application struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *database.ConnectionPool
}
