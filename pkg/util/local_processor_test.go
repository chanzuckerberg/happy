package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarzip(t *testing.T) {
	r := require.New(t)

	processor := NewLocalProcessor()
	curDir, err := os.Getwd()
	r.NoError(err)
	tempFile, err := os.CreateTemp(curDir, "happy_tfe.*.tar.gz")
	r.NoError(err)
	defer os.Remove(tempFile.Name())
	err = processor.Tarzip(".", tempFile)
	r.NoError(err)
	err = processor.Tarzip("/nope", tempFile)
	r.Error(err)
}
