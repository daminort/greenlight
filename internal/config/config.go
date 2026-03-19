package config

import "flag"

type Limiter struct {
	RPS     int
	Burst   int
	Enabled bool
}

type Config struct {
	Port    int
	Env     string
	Version string
	Limiter *Limiter
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

	return &cfg
}
