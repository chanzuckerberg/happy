package util

import "os/exec"

type Executor interface {
	LookPath(file string) (string, error)
	Run(command *exec.Cmd) error
	Output(command *exec.Cmd) ([]byte, error)
}

type DefaultExecutor struct{}

func (e DefaultExecutor) Run(command *exec.Cmd) error {
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
