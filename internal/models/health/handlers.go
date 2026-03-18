package health

import (
	"net/http"

	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/envelopes"
	"greenlight.damian.net/internal/errorsManager"
	"greenlight.damian.net/internal/payloads"
)

type Handlers struct {
	Config       *config.Config
	ErrorManager *errorsManager.ErrorsManager
}

func NewHandlers(cfg *config.Config, em *errorsManager.ErrorsManager) *Handlers {
	return &Handlers{
		Config:       cfg,
		ErrorManager: em,
	}
}

func (h *Handlers) Check(w http.ResponseWriter, r *http.Request) {
	data := Summary{
		Status:      "available",
		Environment: h.Config.Env,
		Version:     h.Config.Version,
	}

	envelope := envelopes.New("summary", data)

	err := payloads.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.ErrorManager.ServerErrorResponse(w, r, err)
	}
}
