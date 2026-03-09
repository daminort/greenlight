package main

import (
	"log/slog"

	"greenlight.damian.net/cmd/api/config"
	"greenlight.damian.net/internal/models"
)

type Application struct {
	Config *config.Config
	Logger *slog.Logger
	Models *models.Models
}
