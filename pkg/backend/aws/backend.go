package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
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
	awsProfile string

	// aws settion: provided or inferred
	awsSession *session.Session

	// aws clients: provided or inferred
	ec2client     ec2iface.EC2API
	ecrclient     ecriface.ECRAPI
	ecsclient     ecsiface.ECSAPI
	logsclient    cloudwatchlogsiface.CloudWatchLogsAPI
	secretsclient secretsmanageriface.SecretsManagerAPI
	ssmclient     ssmiface.SSMAPI
	stsclient     stsiface.STSAPI

	// integration secret: provided or inferred
	integrationSecret *config.IntegrationSecret
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
	if b.awsSession == nil {
		opts := session.Options{
			Profile: b.awsProfile,
			Config: aws.Config{
				Region:     b.awsRegion,
				MaxRetries: aws.Int(2),
			},
			SharedConfigState: session.SharedConfigEnable,
		}

		b.awsSession = session.Must(session.NewSessionWithOptions(opts))
	}

	// Create AWS Clients if we don't have them
	if b.stsclient == nil {
		b.stsclient = sts.New(b.awsSession)
	}

	if b.logsclient == nil {
		b.logsclient = cloudwatchlogs.New(b.awsSession)
	}

	if b.ssmclient == nil {
		b.ssmclient = ssm.New(b.awsSession)
	}

	if b.ecsclient == nil {
		b.ecsclient = ecs.New(b.awsSession)
	}

	if b.ec2client == nil {
		b.ec2client = ec2.New(b.awsSession)
	}

	if b.secretsclient == nil {
		b.secretsclient = secretsmanager.New(b.awsSession)
	}

	if b.ecrclient == nil {
		b.ecrclient = ecr.New(b.awsSession)
	}

	userName, err := b.GetUserName(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve identity info, does aws profile [%s] exist?", b.awsProfile)
	}
	log.Printf("User identity confirmed: %s\n", userName)

	// other inferred or set fields
	if b.integrationSecret == nil {
		integrationSecret, err := b.getIntegrationSecret(ctx, happyConfig.GetSecretArn())
		if err != nil {
			return nil, err
		}
		b.integrationSecret = integrationSecret
	}

	// Create a combined, instantiated config
	b.instantiatedConfig = &instantiatedConfig{
		HappyConfig:       *happyConfig,
		IntegrationSecret: *b.integrationSecret,
	}

	return b, nil
}

func (b *Backend) GetECSClient() ecsiface.ECSAPI {
	return b.ecsclient
}

func (b *Backend) GetEC2Client() ec2iface.EC2API {
	return b.ec2client
}

func (b *Backend) GetECRClient() ecriface.ECRAPI {
	return b.ecrclient
}

func (b *Backend) Conf() *instantiatedConfig {
	return b.instantiatedConfig
}
