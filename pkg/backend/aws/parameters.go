package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TODO: we should already have the path prefix here
func (ab *Backend) GetParam(ctx context.Context, path string) (string, error) {
	logrus.Debugf("reading aws ssm parameter at %s", path)

	out, err := ab.ssmclient.GetParameter(
		ctx,
		&ssm.GetParameterInput{Name: aws.String(path)},
	)
	if err != nil {
		return "", errors.Wrap(err, "could not get parameter")
	}

	return *out.Parameter.Value, nil
}

func (ab *Backend) WriteParam(
	ctx context.Context,
	name string,
	val string,
) error {
	_, err := ab.ssmclient.PutParameter(ctx, &ssm.PutParameterInput{
		Overwrite: true,
		Name:      &name,
		Value:     &val,
	})
	return errors.Wrapf(err, "could not write parameter to %s", name)
}
