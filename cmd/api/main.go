package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	bgManager "greenlight.damian.net/internal/bg_manager"
	"greenlight.damian.net/internal/config"
	"greenlight.damian.net/internal/database"
	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/mailer"
	"greenlight.damian.net/internal/middlewares"
	"greenlight.damian.net/internal/models/health"
	"greenlight.damian.net/internal/models/movies"
	"greenlight.damian.net/internal/models/users"
)

const version = "1.0.0"

func main() {
	// environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// config
	cfg := config.New()

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// error manager
	errorMngr := errorsManager.New(logger)

	// background manager
	bgMngr := bgManager.New(logger)

	// mailer
	mail, err := mailer.New(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.Sender,
	)
	if err != nil {
		logger.Error("Unable to create mailer", "error", err.Error())
		return
	}

	// database
	db, err := database.New()
	if err != nil {
		logger.Error("Unable to connect to database", "error", err.Error())
		return
	}
	defer db.DB.Close()
	logger.Info("database connection pool established")

	// models
	mvRepo := movies.NewRepository(db.DB)
	mvService := movies.NewService(mvRepo)

	usRepo := users.NewRepository(db.DB)
	usService := users.NewService(usRepo)

	// application
	app := &Application{
		Config:       cfg,
		Logger:       logger,
		ErrorManager: errorMngr,
		BgManager:    bgMngr,
		Mailer:       mail,
		Middlewares:  middlewares.New(cfg, errorMngr),
		Movies:       movies.NewHandlers(mvService, errorMngr),
		Users:        users.NewHandlers(usService, errorMngr, mail, bgMngr),
		Health:       health.NewHandlers(cfg, errorMngr),
	}

	err = app.Serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
