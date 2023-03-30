package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	_  = iota //ignore first value by assigning to blank identifier
	kb = 1 << (10 * iota)
	mb
	// GB
)

type Config struct {
	HTTPHost             string `env:"HTTP_HOST" envDefault:"127.0.0.1"`
	HTTPPort             int    `env:"HTTP_PORT" envDefault:"8080"`
	HTTPReqBodySizeLimit int    `env:"HTTP_REQUEST_BODY_SIZE_LIMIT" envDefault:"5"`

	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,notEmpty"`
	DBPassword string `env:"DB_PASSWORD,notEmpty"`
	DBName     string `env:"DB_NAME,notEmpty"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`

	LoggerLevel              string `env:"LOGGER_LEVEL" envDefault:"info"`
	LoggerOutput             string `env:"LOGGER_OUTPUT" envDefault:"stdout"`
	LoggerRotateMaxSize      int    `env:"LOGGER_ROTATE_MAX_SIZE" envDefault:"100"`
	LoggerRotateMaxBackups   int    `env:"LOGGER_ROTATE_MAX_BACKUPS"`
	LoggerRotateWithCompress bool   `env:"LOGGER_ROTATE_WITH_COMPRESS" envDefault:"false"`
	LoggerRotateMaxAge       int    `env:"LOGGER_ROTATE_MAX_AGE"`
}

func Load() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	config.HTTPReqBodySizeLimit *= mb

	return &config, nil
}
