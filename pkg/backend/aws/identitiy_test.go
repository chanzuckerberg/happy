package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUsernameFromCI(t *testing.T) {
	r := require.New(t)

	actor := "FOOBARBAZ"
	t.Setenv("CI", "true")
	t.Setenv("GITHUB_ACTOR", actor)

	b := &Backend{}

	u, err := b.GetUserName(context.Background())
	r.NoError(err)
	r.Equal(actor, u)
}
