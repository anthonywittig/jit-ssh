package sshremoteportforward

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/anthonywittig/jit-ssh/internal/config"
)

type SSHRemotePortForwarder struct {
	cmd                 *exec.Cmd
	commandRunningSince time.Time // We're not being thread safe with this...
}

type Output struct {
	level string
}

func (s *SSHRemotePortForwarder) Start(c config.Config) error {
	if s.Running() {
		log.Println("shouldn't start another process while the first is running; trying to kill it")
		if err := s.Stop(); err != nil {
			log.Printf("error killing running process, not sure what to do, probably best to continue: %s", err.Error())
			s.commandRunningSince = time.Time{}
		}
	}

	log.Println("starting ssh remote port forward")

	go func() {
		s.commandRunningSince = time.Now()

		// The sleep seems to help it not exit immediately (which you don't need to do when running the command manually).
		// Next time you're in here, consider having the remote command be something that echos the uptime every so often.
		command := fmt.Sprintf(`ssh -i %s -o "StrictHostKeyChecking no" -R %d:localhost:22 %s 'sleep 1h'`, c.Local.SSH.PathToKey, c.Remote.PortToOpen, c.Remote.ConnectionString)
		log.Printf("going to execute: %s", command)

		s.cmd = exec.Command("bash", "-s")
		s.cmd.Stdin = strings.NewReader(command)

		s.cmd.Stdout = Output{level: "INFO"}
		s.cmd.Stderr = Output{level: "ERROR"}

		if err := s.cmd.Run(); err != nil {
			log.Printf("error running process (might not be a bad thing): %s", err.Error())
		}
		log.Printf("command exited")

		s.commandRunningSince = time.Time{}
	}()

	return nil
}

func (s *SSHRemotePortForwarder) Running() bool {
	return s.commandRunningSince != time.Time{}
}

func (s *SSHRemotePortForwarder) Stop() error {
	if err := s.cmd.Process.Kill(); err != nil {
		return errors.New(fmt.Sprintf("error killing running process, not sure what to do, probably best to continue: %s", err.Error()))
	}
	s.commandRunningSince = time.Time{}
	return nil
}

func (o Output) Write(b []byte) (int, error) {
	log.Printf("CMD OUTPUT %s - %s", o.level, string(b))
	return len(b), nil
}
