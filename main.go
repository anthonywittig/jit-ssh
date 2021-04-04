package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"

	"github.com/anthonywittig/jit-ssh/internal/application"
	"github.com/anthonywittig/jit-ssh/internal/config/s3configurer"
)

//go:embed .env.json
var embededLocalConfig []byte

func main() {
	if err := run(); err != nil {
		log.Fatalf("error running application: %s", err.Error())
	}
}

func run() error {
	ctx := context.Background()

	configurer, err := s3configurer.New(ctx, embededLocalConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting configurer: %s", err.Error()))
	}

	app, err := application.New(application.Context{
		Configurer: configurer,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("error getting application: %s", err.Error()))
	}

	if err := app.Run(ctx); err != nil {
		return errors.New(fmt.Sprintf("error running application: %s", err.Error()))
	}

	return nil
}
