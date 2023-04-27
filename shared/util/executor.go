package util

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Executor interface {
	LookPath(file string) (string, error)
	Run(command *exec.Cmd) error
	Output(command *exec.Cmd) ([]byte, error)
}

type DefaultExecutor struct{}

func (e DefaultExecutor) Run(command *exec.Cmd) error {
	log.Debugf("executing: %s", command.String())
	return command.Run()
}

func (e DefaultExecutor) Output(command *exec.Cmd) ([]byte, error) {
	return command.Output()
}

func (e DefaultExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func NewDefaultExecutor() Executor {
	return DefaultExecutor{}
}

type DummyExecutor struct{}

func (e DummyExecutor) Run(_ *exec.Cmd) error {
	return nil
}

func (e DummyExecutor) Output(_ *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}

func (e DummyExecutor) LookPath(file string) (string, error) {
	return file, nil
}

func NewDummyExecutor() Executor {
	return DummyExecutor{}
}
