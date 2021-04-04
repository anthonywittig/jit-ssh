package application

import (
	"context"
	"errors"
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
	GetConfig() (config.Config, error)
}

type RemotePortForwarder interface {
	Start() error
	Running() bool
}

func New(ac Context) (*Application, error) {
	return &Application{
		c: ac,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context says it's time to quit, ending run loop")
			return nil
		case <-ticker.C:
			nextDelay, err := a.execute()
			if err != nil {
				return errors.New(fmt.Sprintf("error on app execute: %s", err.Error()))
			}
			ticker.Reset(nextDelay)
		}
	}
}

func (a *Application) execute() (time.Duration, error) {
	conf, err := a.c.Configurer.GetConfig()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error getting config: %s", err.Error()))
	}
	log.Printf("-- config: %+v", conf)

	if conf.Remote.IP == "" {
		return 10 * time.Minute, nil
	}

	// Could we ever get in a state where we're "running" forever?
	if a.c.RemotePortForwarder.Running() {
		return 5 * time.Minute, nil
	}

	a.c.RemotePortForwarder.Start()
	return 1 * time.Minute, nil
}
