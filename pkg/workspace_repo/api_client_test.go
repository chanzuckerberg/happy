package workspace_repo

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func TestGetTfeTokenFromEnv(t *testing.T) {
	r := require.New(t)
	uuid, err := uuid.GenerateUUID()
	r.NoError(err)

	t.Setenv("TFE_TOKEN", uuid)

	token, err := GetTfeToken("")
	r.NoError(err)

	r.Equal(uuid, token)
}

func TestGetTFETokenFromFile(t *testing.T) {
	r := require.New(t)

	// setup things
	token, err := uuid.GenerateUUID()
	r.NoError(err)

	hostname, err := uuid.GenerateUUID()
	r.NoError(err)

	newHome, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(newHome)

	pathToTFEConfig := path.Join(newHome, tfrcFileName)

	// create the dir with the config file
	err = os.MkdirAll(path.Dir(pathToTFEConfig), 0777)
	r.NoError(err)

	// create the config file
	configContents := &tfeConfig{
		Credentials: map[string]tfeCredential{
			hostname: {
				Token: token,
			},
		},
	}

	b, err := json.Marshal(configContents)
	r.NoError(err)

	err = os.WriteFile(pathToTFEConfig, b, 0644)
	r.NoError(err)

	// now try to get the token
	// HACK trick the test by overrideing our "HOME"
	t.Setenv("HOME", newHome)

	gotToken, err := GetTfeToken(hostname)
	r.NoError(err)
	r.Equal(token, gotToken)
}

type tfeCredential struct {
	Token string `json:"token,omitempty"`
}

type tfeConfig struct {
	Credentials map[string]tfeCredential `json:"credentials,omitempty"`
}
