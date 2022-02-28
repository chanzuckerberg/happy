package aws

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestStopError(t *testing.T) {
	r := require.New(t)

	r.False(isStop(nil))
	r.True(isStop(Stop()))
	r.False(isStop(errors.New("foobar")))
}
