package cmd

import (
	"context"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"otusdelivery/build"
	"otusdelivery/config"
)

func postgresCmd(ctx context.Context, conf config.Config) *cobra.Command {
	command := &cobra.Command{ //nolint:exhaustruct
		Use:   "postgres",
		Short: "run db migrations for postgres",
		RunE: func(cmd *cobra.Command, _ []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	command.AddCommand(up(ctx, conf, postgres))

	return command
}

func postgres(ctx context.Context, conf config.Config) (*migrate.Migrate, error) {
	b := build.New(ctx, conf)

	//nolint:wrapcheck
	return b.PostgresMigration()
}

type migrationConstructFn func(context.Context, config.Config) (*migrate.Migrate, error)

func up(ctx context.Context, conf config.Config, constructFn migrationConstructFn) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "up",
		Short: "up migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			m, err := constructFn(ctx, conf)
			if err != nil {
				return errors.Wrap(err, "construct migration")
			}

			err = m.Up()
			if err != nil {
				if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
					return nil
				}

				return errors.Wrap(err, "up migrations")
			}

			return nil
		},
	}
}
