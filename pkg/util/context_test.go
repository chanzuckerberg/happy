package util

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCI(t *testing.T) {
	r := require.New(t)
	ctx, err := BuildContext(context.Background())
	r.NoError(err)

	type c struct {
		ci       *string
		expected bool
	}

	isCI := "true"
	notCI := "false"

	cs := []c{
		{&isCI, true},
		{&notCI, false},
	}

	for idx, c := range cs {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			r := require.New(t)
			if c.ci != nil {
				t.Setenv("CI", *c.ci)
			}

			r.Equal(c.expected, IsCI(ctx))
		})
	}
}
