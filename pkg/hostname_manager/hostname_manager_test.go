package hostname_manager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHostnameManager(t *testing.T) {
	r := require.New(t)

	hostNameManager := NewHostNameManager("./hosts", nil)

	r.NotNil(hostNameManager)
	_, err := hostNameManager.getHostsFileConfig()
	r.NoError(err)

	borders, err := hostNameManager.getFileBorders()
	r.NoError(err)

	config := hostNameManager.generateConfig(borders, []string{"container"})
	r.True(len(config) == 3)
	r.Equal("127.0.0.1\tcontainer", config[1])

	err = hostNameManager.Install()
	r.NoError(err)
	err = hostNameManager.UnInstall()
	r.NoError(err)
}
