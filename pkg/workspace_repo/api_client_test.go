package workspace_repo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWorkspaceRepoErrorNoTFEToken(t *testing.T) {
	r := require.New(t)
	_, err := NewWorkspaceRepo("foo", "bar")
	r.True(err == nil || strings.Contains(err.Error(), "please set env var TFE_TOKEN"))
}
