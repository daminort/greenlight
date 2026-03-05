package main

import (
	"net/http"

	"greenlight.damian.net/internal/utils"
)

type Summary struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := Summary{
		Status:      "available",
		Environment: app.Config.Env,
		Version:     version,
	}

	envelope := utils.NewEnvelope("summary", data)

	err := utils.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
