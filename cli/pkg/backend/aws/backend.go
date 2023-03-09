package aws

import (
	"context"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	compute "github.com/chanzuckerberg/happy/cli/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	awsApiCallMaxRetries   = 100
	awsApiCallBackoffDelay = time.Second * 5
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

	// aws config: provided or inferred
	awsConfig *aws.Config

	// aws clients: provided or inferred
	dynamodbclient              interfaces.DynamoDB
	ec2client                   interfaces.EC2API
	ecrclient                   interfaces.ECRAPI
	ecsclient                   interfaces.ECSAPI
	eksclient                   interfaces.EKSAPI
	secretsclient               interfaces.SecretsManagerAPI
	ssmclient                   interfaces.SSMAPI
	stsclient                   interfaces.STSAPI
	stspresignclient            interfaces.STSPresignAPI
	taskStoppedWaiter           interfaces.ECSTaskStoppedWaiterAPI
	k8sClientCreator            kube.K8sClientCreator
	cwlGetLogEventsAPIClient    interfaces.GetLogEventsAPIClient
	cwlFilterLogEventsAPIClient interfaces.FilterLogEventsAPIClient

	// integration secret: provided or inferred
	integrationSecret    *config.IntegrationSecret
	integrationSecretArn *string

	// cached
	username *string

	executor       util.Executor
	ComputeBackend compute.ComputeBackend
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
		executor:   util.NewDefaultExecutor(),
	}

	b.k8sClientCreator = func(config *rest.Config) (kubernetes.Interface, error) {
		return kubernetes.NewForConfig(config)
	}

	// set optional parameters
	for _, opt := range opts {
		opt(b)
	}

	// Create an AWS session if we don't have one
	if b.awsConfig == nil {
		options := []func(*configv2.LoadOptions) error{configv2.WithRegion(*b.awsRegion),
			configv2.WithRetryer(func() aws.Retryer {
				// Unless specified, we run into ThrottlingException when repeating calls, when following logs or waiting on a condition.
				return retry.AddWithMaxBackoffDelay(retry.AddWithMaxAttempts(retry.NewStandard(), awsApiCallMaxRetries), awsApiCallBackoffDelay)
			})}

		if b.awsProfile != nil && *b.awsProfile != "" {
			options = append(options, configv2.WithSharedConfigProfile(*b.awsProfile))
		}

		if util.IsLocalstackMode() {
			options = append(options, configv2.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               util.GetLocalstackEndpoint(),
						SigningRegion:     region,
						HostnameImmutable: true,
					}, nil
				})))
		}

		conf, err := configv2.LoadDefaultConfig(ctx, options...)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create an aws session")
		}

		b.awsConfig = &conf
	}

	// Create AWS Clients if we don't have them
	if b.stsclient == nil {
		sc := sts.NewFromConfig(*b.awsConfig)
		b.stsclient = sc
		b.stspresignclient = sts.NewPresignClient(sc)
	}

	if b.cwlGetLogEventsAPIClient == nil {
		b.cwlGetLogEventsAPIClient = cwlv2.NewFromConfig(*b.awsConfig)
	}

	if b.cwlFilterLogEventsAPIClient == nil {
		b.cwlFilterLogEventsAPIClient = cwlv2.NewFromConfig(*b.awsConfig)
	}

	if b.ssmclient == nil {
		b.ssmclient = ssm.NewFromConfig(*b.awsConfig)
	}

	if b.ecsclient == nil {
		b.ecsclient = ecs.NewFromConfig(*b.awsConfig)
		b.taskStoppedWaiter = ecs.NewTasksStoppedWaiter(b.ecsclient)
	}

	if b.eksclient == nil {
		b.eksclient = eks.NewFromConfig(*b.awsConfig)
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

	if b.dynamodbclient == nil {
		b.dynamodbclient = dynamodb.NewFromConfig(*b.awsConfig)
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

	b.ComputeBackend, err = b.getComputeBackend(ctx, happyConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to k8s backend")
	}

	// other inferred or set fields
	if b.integrationSecret == nil {
		integrationSecret, integrationSecretArn, err := b.ComputeBackend.GetIntegrationSecret(ctx)
		if err != nil {
			return nil, err
		}
		if integrationSecret.Tfe == nil {
			return nil, errors.New("TFE configuration is not present in the integration secret")
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

func (b Backend) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	return b.awsConfig.Credentials.Retrieve(ctx)
}

func (b *Backend) getComputeBackend(ctx context.Context, happyConfig *config.HappyConfig) (compute.ComputeBackend, error) {
	var computeBackend compute.ComputeBackend
	var err error
	if happyConfig.TaskLaunchType() == util.LaunchTypeK8S {
		computeBackend, err = NewK8SComputeBackend(ctx, *happyConfig.K8SConfig(), b, b.k8sClientCreator)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to connect to k8s backend")
		}
	} else {
		computeBackend, err = NewECSComputeBackend(ctx, happyConfig, b)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to connect to ecs backend")
		}
	}
	return computeBackend, nil
}

func (b *Backend) GetDynamoDBClient() dynamolock.DynamoDBClient {
	return b.dynamodbclient
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

func (b *Backend) GetLogEventsAPIClient() interfaces.GetLogEventsAPIClient {
	return b.cwlGetLogEventsAPIClient
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

func (b *Backend) PrintLogs(ctx context.Context, stackName string, serviceName string, opts ...util.PrintOption) error {
	return b.ComputeBackend.PrintLogs(ctx, stackName, serviceName, opts...)
}

func (b *Backend) RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error {
	return b.ComputeBackend.RunTask(ctx, taskDefArn, launchType)
}

func (b *Backend) Shell(ctx context.Context, stackName string, service string) error {
	return b.ComputeBackend.Shell(ctx, stackName, service)
}

func (b *Backend) GetEvents(ctx context.Context, stackName string, services []string) error {
	return b.ComputeBackend.GetEvents(ctx, stackName, services)
}

func (b *Backend) Describe(ctx context.Context, stackName string, serviceName string) (compute.StackServiceDescription, error) {
	return b.ComputeBackend.Describe(ctx, stackName, serviceName)
}
