package main

import (
	"context"
	"os"

	"otusdelivery/cmd"
	"otusdelivery/config"

	"github.com/rs/zerolog"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	ctx := logger.WithContext(context.Background())

	exitCode := 0

	logger.Info().Msg("application is launching")

	err = cmd.Run(ctx, conf)
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
