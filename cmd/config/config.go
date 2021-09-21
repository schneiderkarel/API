package config

import "app/internal/env"

const (
	envKeyHttpServerPort = "HTTP_SERVER_PORT"
	envKeyPostgresDSN    = "POSTGRES_DSN"
)

type Config struct {
	HttpServerPort int
	PostgresDSN    string
}

func Init() *Config {
	return &Config{
		HttpServerPort: env.MustInt(env.Port(envKeyHttpServerPort, true, "8080")),
		PostgresDSN:    env.MustString(env.String(envKeyPostgresDSN, true, "")),
	}
}
