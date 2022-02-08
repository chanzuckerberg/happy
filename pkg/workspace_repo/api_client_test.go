package workspace_repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWorkspaceRepoErrorNoTFEToken(t *testing.T) {
	r := require.New(t)

	_, err := NewWorkspaceRepo("foo", "bar")
	r.True(err == nil || err.Error() == "Please set env var TFE_TOKEN")
}
