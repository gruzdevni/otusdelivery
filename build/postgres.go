//nolint:revive
package build

import (
	"context"
	"fmt"
	"net"

	_ "github.com/jackc/pgx/stdlib" // driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const dsnTemplate = "postgres://%s:%s@%s/%s?sslmode=disable"

func (b *Builder) postgresClient(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to postgres")
	}

	db.SetMaxOpenConns(b.config.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(b.config.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(b.config.Postgres.ConnMaxLifetime)

	b.shutdown.add(func(_ context.Context) error {
		if err = db.Close(); err != nil {
			return errors.Wrap(err, "close db connection")
		}

		return nil
	})

	return db, nil
}

func (b *Builder) PostgresClient() (*sqlx.DB, error) {
	return b.postgresClient(b.PostgresDSN())
}

func (b *Builder) PostgresDSN() string {
	return fmt.Sprintf(dsnTemplate,
		b.config.Postgres.DsnUser,
		b.config.Postgres.DsnPassword,
		net.JoinHostPort(b.config.Postgres.DsnHost, b.config.Postgres.DsnPort),
		b.config.Postgres.DsnDBName,
	)
}
