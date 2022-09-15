package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
)

func (b *Backend) getIntegrationSecret(ctx context.Context, happyConfig *config.HappyConfig) (*config.IntegrationSecret, *string, error) {
	secretId := happyConfig.GetSecretId()
	out, err := b.secretsclient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get integration secret at %s", secretId)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), secret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
	}
	return secret, out.ARN, nil
}
