package config

import (
	"flag"
	"os"
	"strconv"
)

type Limiter struct {
	RPS     int
	Burst   int
	Enabled bool
}

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type Config struct {
	Port    int
	Env     string
	Version string
	Limiter *Limiter
	SMTP    *SMTP
}

func New() *Config {
	cfg := Config{
		Version: "0.0.1",
	}

	limiter := &Limiter{}

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|stage|prod)")

	flag.IntVar(&limiter.RPS, "limiter-rps", 2, "Rate limit: maximum requests per second")
	flag.IntVar(&limiter.Burst, "limiter-burst", 4, "Rate limit: maximum burst")
	flag.BoolVar(&limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	cfg.Limiter = limiter

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpSender := os.Getenv("SMTP_SENDER")

	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		smtpPortInt = 2525
	}

	smtp := &SMTP{
		Host:     smtpHost,
		Port:     smtpPortInt,
		Username: smtpUsername,
		Password: smtpPassword,
		Sender:   smtpSender,
	}

	cfg.SMTP = smtp

	return &cfg
}
