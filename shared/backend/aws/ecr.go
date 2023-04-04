package aws

import (
	"context"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/pkg/errors"
)

type ECRAuthorizationToken struct {
	Username      string
	Password      string
	ProxyEndpoint string
}

// NOTE: we just need one token to access al ECRs this principal has access to
func (b *Backend) ECRGetAuthorizationToken(ctx context.Context) (*ECRAuthorizationToken, error) {
	encodedTokens, err := b.ecrclient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, errors.Wrap(err, "could not get ECR authorizaiton token")
	}

	// NOTE: because registryIDs is deprecated, we assume there is only one token being generated that can
	//       be used for all our registries.
	authData := encodedTokens.AuthorizationData[0]
	decodedToken, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
	if err != nil {
		return nil, errors.Wrap(err, "could not base64 decode ECR authorization token")
	}

	split := strings.Split(string(decodedToken), ":")
	username := split[0]
	password := split[1]

	// ProxyEndpoint is the registry URL to use for the authorization token in a docker login command
	// We need to transform this to be compatible with a docker login command
	endpoint, err := url.Parse(*authData.ProxyEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse docker registry URL")
	}

	return &ECRAuthorizationToken{
		Username:      username,
		Password:      password,
		ProxyEndpoint: endpoint.Host,
	}, nil
}
