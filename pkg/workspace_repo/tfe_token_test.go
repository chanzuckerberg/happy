package workspace_repo

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTfeToken(t *testing.T) {
	req := require.New(t)

	dir, err := os.Getwd()
	req.NoError(err)
	t.Setenv("HOME", path.Join(dir, "testdata"))
	token, err := GetTfeToken("https://www.example.com")
	req.NoError(err)
	req.Equal("aaa.bbb.ccc", token)
	t.Setenv("HOME", path.Join(dir, "testdata_nope"))
	_, err = GetTfeToken("https://www.example.com")
	req.Error(err)
}
