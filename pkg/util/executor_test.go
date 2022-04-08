package util

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecutor(t *testing.T) {
	r := require.New(t)
	execPath, err := exec.LookPath("pwd")
	r.NoError(err)
	cmd := exec.CommandContext(context.Background(), execPath)

	executor := NewDummyExecutor()
	err = executor.Run(cmd)
	r.NoError(err)
	cmd = exec.CommandContext(context.Background(), execPath)
	_, err = executor.Output(cmd)
	r.NoError(err)
	found, err := executor.LookPath("foobar")
	r.NoError(err)
	r.NotEmpty(found)

	executor = NewDefaultExecutor()
	err = executor.Run(cmd)
	r.NoError(err)
	cmd = exec.CommandContext(context.Background(), execPath)
	_, err = executor.Output(cmd)
	r.NoError(err)
	_, err = executor.LookPath("foobar")
	r.Error(err)
}
