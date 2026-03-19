package middlewares

import (
	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/errors_manager"
)

type Middlewares struct {
	Config       *config.Config
	ErrorManager *errorsManager.ErrorsManager
}

func New(cfg *config.Config, em *errorsManager.ErrorsManager) *Middlewares {
	return &Middlewares{
		Config:       cfg,
		ErrorManager: em,
	}
}
