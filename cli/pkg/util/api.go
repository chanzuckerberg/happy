package util

import (
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/shared/client"
)

func MakeApiClient(happyConfig *config.HappyConfig) *client.HappyClient {
	return client.NewHappyClient(
		"happy",
		GetVersion().Version,
		happyConfig.GetHappyApiConfig().BaseUrl,
		happyConfig.GetHappyApiConfig().OidcClientID,
		happyConfig.GetHappyApiConfig().OidcClientID,
	)
}
