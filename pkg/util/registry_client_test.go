package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryClient(t *testing.T) {
	r := require.New(t)

	client := NewDefaultRegistryClient()
	err := client.Login("foo", "bar", "http://registry.com")
	r.Error(err)

	client = NewDummyRegistryClient()
	err = client.Login("foo", "bar", "http://registry.com")
	r.NoError(err)
}
