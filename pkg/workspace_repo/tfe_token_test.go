package workspace_repo

import (
	"os"
	"path"
	"testing"

	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestTfeToken(t *testing.T) {
	r := require.New(t)
	os.Unsetenv("TFE_TOKEN")
	dir, err := os.Getwd()
	r.NoError(err)
	os.Setenv("HOME", path.Join(dir, "mock_data"))
	token, err := GetTfeToken("https://www.tfe.com", util.NewDummyExecutor())
	r.NoError(err)
	r.Equal("aaa.bbb.ccc", token)
	os.Setenv("HOME", path.Join(dir, "mock_data_nope"))
	_, err = GetTfeToken("https://www.tfe.com", util.NewDummyExecutor())
	r.Error(err)
}
