// Package config
package config

import (
	"medods-jwt/internal/trasnport/rest"
	"medods-jwt/pkg/db/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config
type Config struct {
	postgres.PostgresConfig
	rest.ServerConfig
	Debug     bool   `env:"DEBUG" env-default:"false"`
	SecretKey string `env:"SHA512KEY" env-default:"secrets"`
}

// New
func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
