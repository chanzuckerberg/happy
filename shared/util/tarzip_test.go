package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarzip(t *testing.T) {
	r := require.New(t)

	curDir, err := os.Getwd()
	r.NoError(err)
	tempFile, err := os.CreateTemp(curDir, "happy_tfe.*.tar.gz")
	r.NoError(err)
	defer os.Remove(tempFile.Name())
	err = TarDir(".", tempFile)
	r.NoError(err)
	err = TarDir("/nope", tempFile)
	r.Error(err)
}
