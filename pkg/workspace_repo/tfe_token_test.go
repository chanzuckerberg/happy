package workspace_repo

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTfeToken(t *testing.T) {
	req := require.New(t)

	os.Unsetenv("TFE_TOKEN")
	dir, err := os.Getwd()
	req.NoError(err)
	os.Setenv("HOME", path.Join(dir, "testdata"))
	token, err := GetTfeToken(context.Background(), "https://www.tfe.com")
	req.NoError(err)
	req.Equal("aaa.bbb.ccc", token)
	os.Setenv("HOME", path.Join(dir, "testdata_nope"))
	_, err = GetTfeToken(context.Background(), "https://www.tfe.com")
	req.Error(err)
}
