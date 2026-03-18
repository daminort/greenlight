package config

import "flag"

type Config struct {
	Port    int
	Env     string
	Version string
}

func New() *Config {
	cfg := Config{
		Version: "0.0.1",
	}

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	return &cfg
}
