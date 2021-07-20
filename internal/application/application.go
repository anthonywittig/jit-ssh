package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anthonywittig/jit-ssh/internal/config"
)

type Application struct {
	c Context
}

type Context struct {
	Configurer          Configurer
	RemotePortForwarder RemotePortForwarder
}

type Configurer interface {
	GetConfig(context.Context) (config.Config, error)
}

type RemotePortForwarder interface {
	Running() bool
	Start(config.Config) error
	Stop() error
}

func New(ac Context) (*Application, error) {
	return &Application{
		c: ac,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context says it's time to quit, ending run loop")
			return nil
		case <-ticker.C:
			nextDelay, err := a.execute(ctx)
			if err != nil {
				return fmt.Errorf("error on app execute: %s", err.Error())
			}
			log.Printf("sleeping for %.1f minutes", nextDelay.Minutes())
			ticker.Reset(nextDelay)
		}
	}
}

func (a *Application) execute(ctx context.Context) (time.Duration, error) {
	conf, err := a.c.Configurer.GetConfig(ctx)
	if err != nil {
		return -1, fmt.Errorf("error getting config: %s", err.Error())
	}

	if conf.Remote.ConnectionString == "" {
		if err := a.c.RemotePortForwarder.Stop(); err != nil {
			return -1, fmt.Errorf("error trying to stop port forwarder: %s", err.Error())
		}
		return 5 * time.Minute, nil
	}

	// Could we ever get in a state where we're "running" forever?
	if a.c.RemotePortForwarder.Running() {
		return 1 * time.Minute, nil
	}

	if err := a.c.RemotePortForwarder.Start(conf); err != nil {
		return -1, fmt.Errorf("error trying to start port forwarder: %s", err.Error())
	}

	return 1 * time.Minute, nil
}
