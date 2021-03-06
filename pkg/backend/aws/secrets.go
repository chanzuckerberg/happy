package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
)

func (b *Backend) getIntegrationSecret(ctx context.Context, secretARN string) (*config.IntegrationSecret, *string, error) {
	out, err := b.secretsclient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretARN,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get integration secret at %s", secretARN)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), secret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
	}
	return secret, out.ARN, nil
}
