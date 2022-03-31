package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type instantiatedConfig struct {
	config.HappyConfig
	config.IntegrationSecret
}

type Backend struct {
	// requiredInputs
	instantiatedConfig *instantiatedConfig

	// optional inputs
	awsRegion  *string
	awsProfile *string

	awsAccountID *string

	// aws settion: provided or inferred
	awsConfig *aws.Config

	// aws clients: provided or inferred
	ec2client         interfaces.EC2API
	ecrclient         interfaces.ECRAPI
	ecsclient         interfaces.ECSAPI
	secretsclient     interfaces.SecretsManagerAPI
	ssmclient         interfaces.SSMAPI
	stsclient         interfaces.STSAPI
	taskRunningWaiter interfaces.ECSTaskRunningWaiterAPI
	taskStoppedWaiter interfaces.ECSTaskStoppedWaiterAPI

	// aws v2 clients: provided or inferred
	cwlGetLogEventsAPIClient interfaces.GetLogEventsAPIClient

	// integration secret: provided or inferred
	integrationSecret    *config.IntegrationSecret
	integrationSecretArn *string

	// cached
	username *string
}

// New returns a new AWS backend
func NewAWSBackend(
	ctx context.Context,
	happyConfig *config.HappyConfig,
	opts ...AWSBackendOption) (*Backend, error) {
	// Set defaults
	b := &Backend{
		awsRegion:  aws.String("us-west-2"),
		awsProfile: happyConfig.AwsProfile(),
	}

	// set optional parameters
	for _, opt := range opts {
		opt(b)
	}

	// Create an AWS session if we don't have one
	if b.awsConfig == nil {
		conf, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(*b.awsRegion), configv2.WithSharedConfigProfile(*b.awsProfile), configv2.WithRetryMaxAttempts(2))
		if err != nil {
			return nil, errors.Wrap(err, "unable to create an aws session")
		}

		b.awsConfig = &conf
	}

	// NOTE: we also create an aws sdk V2 client as we migrate over
	v2Opts := []func(*configv2.LoadOptions) error{
		configv2.WithRegion(*b.awsRegion),
	}
	if b.awsProfile != nil && *b.awsProfile != "" {
		v2Opts = append(v2Opts, configv2.WithSharedConfigProfile(*b.awsProfile))
	}
	cfg, err := configv2.LoadDefaultConfig(
		ctx,
		v2Opts...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load aws SDK config")
	}

	// Create AWS Clients if we don't have them
	if b.stsclient == nil {
		b.stsclient = sts.NewFromConfig(*b.awsConfig)
	}

	if b.cwlGetLogEventsAPIClient == nil {
		b.cwlGetLogEventsAPIClient = cwlv2.NewFromConfig(cfg)
	}

	if b.ssmclient == nil {
		b.ssmclient = ssm.NewFromConfig(*b.awsConfig)
	}

	if b.ecsclient == nil {
		b.ecsclient = ecs.NewFromConfig(*b.awsConfig)
		b.taskRunningWaiter = ecs.NewTasksRunningWaiter(b.ecsclient)
		b.taskStoppedWaiter = ecs.NewTasksStoppedWaiter(b.ecsclient)
	}

	if b.ec2client == nil {
		b.ec2client = ec2.NewFromConfig(*b.awsConfig)
	}

	if b.secretsclient == nil {
		b.secretsclient = secretsmanager.NewFromConfig(*b.awsConfig)
	}

	if b.ecrclient == nil {
		b.ecrclient = ecr.NewFromConfig(*b.awsConfig)
	}

	userName, err := b.GetUserName(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve identity info, does aws profile [%s] exist?", *b.awsProfile)
	}
	logrus.Debugf("user identity confirmed: %s\n", userName)

	accountID, err := b.GetAccountID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve aws account id")
	}
	logrus.Debugf("AWS accunt ID confirmed: %s\n", accountID)

	// other inferred or set fields
	if b.integrationSecret == nil {
		integrationSecret, integrationSecretArn, err := b.getIntegrationSecret(ctx, happyConfig.GetSecretArn())
		if err != nil {
			return nil, err
		}
		b.integrationSecret = integrationSecret
		b.integrationSecretArn = integrationSecretArn
	}

	// Create a combined, instantiated config
	b.instantiatedConfig = &instantiatedConfig{
		HappyConfig:       *happyConfig,
		IntegrationSecret: *b.integrationSecret,
	}

	return b, nil
}

func (b *Backend) GetECSClient() interfaces.ECSAPI {
	return b.ecsclient
}

func (b *Backend) GetEC2Client() interfaces.EC2API {
	return b.ec2client
}

func (b *Backend) GetECRClient() interfaces.ECRAPI {
	return b.ecrclient
}

func (b *Backend) Conf() *instantiatedConfig {
	return b.instantiatedConfig
}

func (b *Backend) GetAWSRegion() string {
	return *b.awsRegion
}

func (b *Backend) GetAWSProfile() string {
	return *b.awsProfile
}

func (b *Backend) GetAWSAccountID() string {
	return *b.awsAccountID
}

func (b *Backend) GetIntegrationSecret() *config.IntegrationSecret {
	return b.integrationSecret
}

func (b *Backend) GetIntegrationSecretArn() *string {
	return b.integrationSecretArn
}
