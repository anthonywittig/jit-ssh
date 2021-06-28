package application

import (
	"context"
	"testing"

	"github.com/anthonywittig/jit-ssh/internal/config"
)

func TestApplicationExecute_callsGetConfig(t *testing.T) {
	mConfigurer := &mockConfigurer{
		getConfig: func() (config.Config, error) { return config.Config{}, nil },
	}
	mRPF := &mockRPF{
		stop: func() error { return nil },
	}
	app, err := New(Context{
		Configurer:          mConfigurer,
		RemotePortForwarder: mRPF,
	})
	if err != nil {
		t.Fatalf("error getting app: %s", err.Error())
	}

	if mConfigurer.getConfigCallCount != 0 {
		t.Fatalf("expected GetConfig call count to be 0 but was %d", mConfigurer.getConfigCallCount)
	}

	if _, err := app.execute(context.Background()); err != nil {
		t.Fatalf("error calling execute: %s", err.Error())
	}

	if mConfigurer.getConfigCallCount != 1 {
		t.Fatalf("expected GetConfig call count to be 1 but was %d", mConfigurer.getConfigCallCount)
	}
}

func TestApplicationExecute_hasRemoteIPIsRunning(t *testing.T) {
	mConfigurer := &mockConfigurer{
		getConfig: func() (config.Config, error) {
			return config.Config{
				Remote: config.Remote{
					ConnectionString: "ubuntu@ec2-13-56-168-223.us-west-1.compute.amazonaws.com",
				},
			}, nil
		},
	}
	mRPF := &mockRPF{
		running: func() bool { return true },
	}
	app, err := New(Context{
		Configurer:          mConfigurer,
		RemotePortForwarder: mRPF,
	})
	if err != nil {
		t.Fatalf("error getting app: %s", err.Error())
	}

	if mRPF.startCallCount != 0 || mRPF.runningCallCount != 0 {
		t.Fatalf("unexpected call counts: %d, %d", mRPF.startCallCount, mRPF.runningCallCount)
	}

	if _, err := app.execute(context.Background()); err != nil {
		t.Fatalf("error calling execute: %s", err.Error())
	}

	if mRPF.startCallCount != 0 || mRPF.runningCallCount != 1 {
		t.Fatalf("unexpected call counts: %d, %d", mRPF.startCallCount, mRPF.runningCallCount)
	}
}

func TestApplicationExecute_hasRemoteIPNotRunning(t *testing.T) {
	mConfigurer := &mockConfigurer{
		getConfig: func() (config.Config, error) {
			return config.Config{
				Remote: config.Remote{
					ConnectionString: "ubuntu@ec2-13-56-168-223.us-west-1.compute.amazonaws.com",
				},
			}, nil
		},
	}
	mRPF := &mockRPF{
		start:   func() error { return nil },
		running: func() bool { return false },
	}
	app, err := New(Context{
		Configurer:          mConfigurer,
		RemotePortForwarder: mRPF,
	})
	if err != nil {
		t.Fatalf("error getting app: %s", err.Error())
	}

	if mRPF.startCallCount != 0 || mRPF.runningCallCount != 0 {
		t.Fatalf("unexpected call counts: %d, %d", mRPF.startCallCount, mRPF.runningCallCount)
	}

	if _, err := app.execute(context.Background()); err != nil {
		t.Fatalf("error calling execute: %s", err.Error())
	}

	if mRPF.startCallCount != 1 || mRPF.runningCallCount != 1 {
		t.Fatalf("unexpected call counts: %d, %d", mRPF.startCallCount, mRPF.runningCallCount)
	}
}

type mockConfigurer struct {
	getConfig          func() (config.Config, error)
	getConfigCallCount int
}

func (m *mockConfigurer) GetConfig(ctx context.Context) (config.Config, error) {
	if m.getConfig == nil {
		panic("need to set up the mock")
	}
	m.getConfigCallCount++
	return m.getConfig()
}

type mockRPF struct {
	running          func() bool
	runningCallCount int
	start            func() error
	startCallCount   int
	stop             func() error
	stopCallCount    int
}

func (m *mockRPF) Running() bool {
	if m.running == nil {
		panic("need to set up the mock")
	}
	m.runningCallCount++
	return m.running()
}

func (m *mockRPF) Start(cfg config.Config) error {
	if m.start == nil {
		panic("need to set up the mock")
	}
	m.startCallCount++
	return m.start()
}

func (m *mockRPF) Stop() error {
	if m.stop == nil {
		panic("need to set up the mock")
	}
	m.stopCallCount++
	return m.stop()
}
