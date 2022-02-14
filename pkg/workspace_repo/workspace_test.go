package workspace_repo

import (
	"testing"

	"os"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceRepo(t *testing.T) {
	r := require.New(t)
	os.Setenv("TFE_TOKEN", "token")
	repo, err := NewWorkspaceRepo("https://repo.com", "organization")
	r.NoError(err)
	_, err = repo.getToken("hostname")
	r.NoError(err)
	_, err = repo.getTfc()
	r.NoError(err)
	_, err = repo.Stacks()
	r.NoError(err)
}
