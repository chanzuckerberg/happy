package util

import (
	"context"

	oidc "github.com/chanzuckerberg/go-misc/oidc_cli"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/pkg/errors"
)

type CliTokenProvider struct {
	oidcClientID  string
	oidcIssuerURL string
}

func (t CliTokenProvider) GetToken() (string, error) {
	token, err := oidc.GetToken(context.Background(), t.oidcClientID, t.oidcIssuerURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	return token.IDToken, nil
}

func MakeApiClient(happyConfig *config.HappyConfig) *client.HappyClient {
	tokenProvider := CliTokenProvider{
		oidcClientID:  happyConfig.GetHappyApiConfig().OidcClientID,
		oidcIssuerURL: happyConfig.GetHappyApiConfig().OidcIssuerUrl,
	}
	return client.NewHappyClient(
		"happy",
		GetVersion().Version,
		happyConfig.GetHappyApiConfig().BaseUrl,
		tokenProvider,
	)
}
