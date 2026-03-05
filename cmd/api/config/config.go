package config

import "flag"

type Config struct {
	Port int
	Env  string
}

func New() *Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	return &cfg
}
