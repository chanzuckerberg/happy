package aws

import (
	"context"
	"testing"
	"time"

	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/stretchr/testify/require"
)

func TestGetUsernameFromCI(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())

	actor := "FOOBARBAZ"
	t.Setenv("CI", "true")
	t.Setenv("GITHUB_ACTOR", actor)

	b := &Backend{}

	u, err := b.GetUserName(context.Background())
	r.NoError(err)
	r.Equal(actor, u)

	tag, err := b.GenerateTag(ctx)
	r.NoError(err)
	r.NotEmpty(tag)
}
