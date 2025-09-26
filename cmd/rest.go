package cmd

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"otusdelivery/build"
	"otusdelivery/config"
)

func restCmd(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "rest",
		Short: "start rest server",
		RunE: func(_ *cobra.Command, _ []string) error {
			builder := build.New(ctx, conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			server, err := builder.RestAPIServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build rest api server")
			}

			if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return errors.Wrap(err, "rest api server serve")
			}

			<-ctx.Done()

			return nil
		},
	}
}
