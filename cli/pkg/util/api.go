package util

import (
	"github.com/chanzuckerberg/happy-shared/client"
	"github.com/chanzuckerberg/happy/pkg/config"
)

func MakeApiClient(happyConfig *config.HappyConfig) *client.HappyClient {
	return client.NewHappyClient("happy", GetVersion().Version, happyConfig.GetHappyApiBaseUrl())
}
