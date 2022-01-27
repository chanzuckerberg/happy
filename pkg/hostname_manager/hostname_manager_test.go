package hostname_manager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHostnameManager(t *testing.T) {
	r := require.New(t)

	r.NotNil(NewHostNameManager("foo", nil))
}
