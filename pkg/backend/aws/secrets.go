package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
)

func (b *Backend) getIntegrationSecret(ctx context.Context, secretARN string) (*config.IntegrationSecret, error) {
	out, err := b.secretsclient.GetSecretValueWithContext(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretARN,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not get integration secret at %s", secretARN)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal(out.SecretBinary, secret)
	if err != nil {
		return nil, errors.Wrap(err, "could not json parse integraiton secret")
	}
	return secret, nil
}
