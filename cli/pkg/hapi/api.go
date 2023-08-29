package hapi

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/chanzuckerberg/go-misc/oidc_cli/oidc_impl"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
)

type CliTokenProvider struct {
	oidcClientID  string
	oidcIssuerURL string
}

const HappyOIDCIDTokenEnvVar = "HAPPY_OIDC_ID_TOKEN"

func (t CliTokenProvider) GetToken() (string, error) {
	// if an environment variable is set, use that instead of the CLI
	// this allows us to use the CLI in CI
	if token := os.Getenv(HappyOIDCIDTokenEnvVar); token != "" {
		return token, nil
	}
	token, err := oidc_impl.GetToken(context.Background(), t.oidcClientID, t.oidcIssuerURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	return token.IDToken, nil
}

type AWSCredentialsProviderCLI struct {
	backend *backend.Backend
}

func (c AWSCredentialsProviderCLI) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	return c.backend.GetCredentials(ctx)
}

type APIClientOption func(*client.HappyClient)

func WithBaseURL(s string) APIClientOption {
	return func(c *client.HappyClient) {
		c.APIBaseUrl = s
	}
}

func MakeAPIClient(happyConfig *config.HappyConfig, backend *backend.Backend, opts ...APIClientOption) *client.HappyClient {
	tokenProvider := CliTokenProvider{
		oidcClientID:  happyConfig.GetHappyAPIConfig().OidcClientID,
		oidcIssuerURL: happyConfig.GetHappyAPIConfig().OidcIssuerUrl,
	}
	awsCredsProvider := AWSCredentialsProviderCLI{
		backend: backend,
	}

	happyClient := client.NewHappyClient(
		"happy-cli",
		util.GetVersion().Version,
		happyConfig.GetHappyAPIConfig().BaseUrl,
		tokenProvider,
		awsCredsProvider,
	)

	for _, opt := range opts {
		opt(happyClient)
	}

	return happyClient
}
