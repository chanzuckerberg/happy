package util

import (
	"context"
	"fmt"

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
	fmt.Println("...> getting token")
	fmt.Println("... - client  id", t.oidcClientID)
	fmt.Println("... - client url", t.oidcIssuerURL)
	token, err := oidc.GetToken(context.Background(), t.oidcClientID, t.oidcIssuerURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	fmt.Println("...> token created")

	tokenStr, err := token.Marshal()
	fmt.Println("...> token marshaled")
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal token")
	}
	fmt.Println("...> token str", tokenStr)

	return tokenStr, nil
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
