package util

import (
	"github.com/chanzuckerberg/happy/pkg/cli/config"
	"github.com/chanzuckerberg/happy/shared/client"
)

func MakeApiClient(happyConfig *config.HappyConfig) *client.HappyClient {
	return client.NewHappyClient("happy", GetVersion().Version, happyConfig.GetHappyApiBaseUrl())
}
