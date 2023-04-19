package util

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestStopError(t *testing.T) {
	r := require.New(t)

	r.False(IsStop(nil))
	r.True(IsStop(Stop()))
	r.False(IsStop(errors.New("foobar")))
}
