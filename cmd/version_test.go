package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	r := require.New(t)
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"version"})
	rootCmd.Execute()
	out, err := io.ReadAll(b)
	r.NoError(err)
	r.Equal("version: undefined\ngit_sha: undefined", string(out))
}
