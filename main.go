package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/anthonywittig/jit-ssh/internal/application"
	"github.com/anthonywittig/jit-ssh/internal/config/s3configurer"
	"github.com/anthonywittig/jit-ssh/internal/remoteportforward/sshremoteportforward"
)

func main() {
	for {
		if err := run(); err != nil {
			log.Fatalf("error running application (sleeping for a minute): %s", err.Error())
		}
		time.Sleep(time.Minute)
	}
}

func run() error {
	ctx := context.Background()

	configurer, err := s3configurer.New(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting configurer: %s", err.Error()))
	}

	app, err := application.New(application.Context{
		Configurer:          configurer,
		RemotePortForwarder: &sshremoteportforward.SSHRemotePortForwarder{},
	})
	if err != nil {
		return errors.New(fmt.Sprintf("error getting application: %s", err.Error()))
	}

	if err := app.Run(ctx); err != nil {
		return errors.New(fmt.Sprintf("error running application: %s", err.Error()))
	}

	return nil
}
