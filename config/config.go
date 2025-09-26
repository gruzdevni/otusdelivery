package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	App      App
	HTTP     HTTP
	Postgres Postgres
}

type appEnv string

const (
	appEnvLocal appEnv = "local"
	appEnvDev   appEnv = "dev"
	appEnvProd  appEnv = "prod"
)

type App struct {
	ENV  appEnv `envconfig:"APP_ENV"                default:"local"`
	Name string `envconfig:"APP_NAME"               default:"app"`
}

type HTTP struct {
	Port    int32    `envconfig:"HTTP_PORT" default:"8080"`
	Schemes []string `envconfig:"HTTP_SCHEMES" default:"http"`
}

func Load() (Config, error) {
	cnf := Config{} //nolint:exhaustruct

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cnf, errors.Wrap(err, "read .env file")
	}

	if err := envconfig.Process("", &cnf); err != nil {
		return cnf, errors.Wrap(err, "read environment")
	}

	return cnf, nil
}

func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTP.Port)
}
