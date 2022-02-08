package aws

import (
	"context"
	"encoding/base64"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
)

type ECRAuthorizationToken struct {
	Registry string
	Username string
	Password string
}

func (e *ECRAuthorizationToken) DockerLogin(ctx context.Context) error {
	args := []string{"docker", "login", "--username", e.Username, "--password-stdin", e.Registry}

	docker, err := exec.LookPath("docker")
	if err != nil {
		return errors.Wrap(err, "could not find docker in path")
	}
	cmd := exec.CommandContext(ctx, docker, args...)
	cmd.Stdin = strings.NewReader(e.Password)

	err = cmd.Run()
	return errors.Wrap(err, "registry login failed")
}

func (b *Backend) ECRGetAuthorizationTokens(ctx context.Context, registryIDs []string) ([]ECRAuthorizationToken, error) {
	registryIDPtrs := []*string{}
	for _, rID := range registryIDs {
		registryIDPtrs = append(registryIDPtrs, aws.String(rID))
	}

	encodedTokens, err := b.ecrclient.GetAuthorizationTokenWithContext(ctx, &ecr.GetAuthorizationTokenInput{
		RegistryIds: registryIDPtrs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get ECR authorizaiton token")
	}

	tokens := []ECRAuthorizationToken{}
	for idx, registry := range registryIDs {
		decodedToken, err := base64.StdEncoding.DecodeString(*encodedTokens.AuthorizationData[idx].AuthorizationToken)
		if err != nil {
			return nil, errors.Wrap(err, "could not base64 decode ECR authorization token")
		}

		split := strings.Split(string(decodedToken), ":")
		username := split[0]
		password := split[1]

		tokens = append(tokens, ECRAuthorizationToken{
			Registry: registry,
			Username: username,
			Password: password,
		})
	}

	return tokens, nil
}
