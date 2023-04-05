package hapi

import (
	"context"

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

func (t CliTokenProvider) GetToken() (string, error) {
	token, err := oidc_impl.GetToken(context.Background(), t.oidcClientID, t.oidcIssuerURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	return token.IDToken, nil
}

type AWSCredentialsProviderCLI struct {
	happyConfig *config.HappyConfig
}

func (c AWSCredentialsProviderCLI) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	b, err := backend.NewAWSBackend(ctx, c.happyConfig.GetEnvironmentContext())
	if err != nil {
		return aws.Credentials{}, err
	}
	return b.GetCredentials(ctx)
}

func MakeApiClient(happyConfig *config.HappyConfig) *client.HappyClient {
	tokenProvider := CliTokenProvider{
		oidcClientID:  happyConfig.GetHappyApiConfig().OidcClientID,
		oidcIssuerURL: happyConfig.GetHappyApiConfig().OidcIssuerUrl,
	}
	awsCredsProvider := AWSCredentialsProviderCLI{
		happyConfig: happyConfig,
	}

	return client.NewHappyClient(
		"happy-cli",
		util.GetVersion().Version,
		happyConfig.GetHappyApiConfig().BaseUrl,
		tokenProvider,
		awsCredsProvider,
	)
}
