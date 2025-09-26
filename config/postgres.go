package config

import "time"

type Postgres struct {
	DsnHost     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	DsnPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	DsnDBName   string `envconfig:"POSTGRES_DB_NAME" default:"postgres"`
	DsnUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	DsnPassword string `envconfig:"POSTGRES_PASSWORD" default:""`

	MaxOpenConns    int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns    int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"7"`
	ConnMaxLifetime time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"30m"`
}
