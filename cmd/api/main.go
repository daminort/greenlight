package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"greenlight.damian.net/cmd/api/config"
)

const version = "1.0.0"

func main() {

	// config
	cfg := config.New()

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &Application{
		Config: cfg,
		Logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
