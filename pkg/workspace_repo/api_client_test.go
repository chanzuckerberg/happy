package workspace_repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWorkspaceRepoErrorNoTFEToken(t *testing.T) {
	r := require.New(t)

	_, err := NewWorkspaceRepo("foo", "bar")
	r.Error(err, "Please set env var TFE_TOKEN")
}
