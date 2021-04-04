package sshremoteportforward

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type SSHRemotePortForwarder struct {
	cmd                 *exec.Cmd
	commandRunningSince time.Time
}

func (s *SSHRemotePortForwarder) Start() error {
	if s.Running() {
		log.Println("shouldn't start another process while the first is running; trying to kill it")
		if err := s.Stop(); err != nil {
			log.Printf("error killing running process, not sure what to do, probably best to continue: %s", err.Error())
			s.commandRunningSince = time.Time{}
		}
	}

	log.Println("starting ssh remote port forward")
	log.Println("  (implement me!)")

	go func() {
		s.commandRunningSince = time.Now()

		s.cmd = exec.Command("sleep", "5")
		if err := s.cmd.Run(); err != nil {
			log.Printf("error running process (might not be a bad thing): %s", err.Error())
		}

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
