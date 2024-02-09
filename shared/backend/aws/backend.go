package aws

import (
	"context"
	"fmt"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	compute "github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	awsApiCallMaxRetries   = 100
	awsApiCallBackoffDelay = time.Second * 5
)

type instantiatedConfig struct {
	config.EnvironmentContext
	config.IntegrationSecret
}

type Backend struct {
	// requiredInputs
	instantiatedConfig *instantiatedConfig

	// optional inputs
	awsRegion  *string
	awsProfile *string

	awsAccountID *string

	awsRoleArn *string

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

	executor             util.Executor
	computeBackend       compute.ComputeBackend
	environmentContext   config.EnvironmentContext
	awsConfigLoadOptions []func(*configv2.LoadOptions) error
}

// New returns a new AWS backend
func NewAWSBackend(
	ctx context.Context,
	environmentContext config.EnvironmentContext,
	opts ...AWSBackendOption) (*Backend, error) {
	// Set defaults
	b := &Backend{
		awsRegion:          environmentContext.AWSRegion,
		awsProfile:         environmentContext.AWSProfile,
		executor:           util.NewDefaultExecutor(),
		environmentContext: environmentContext,
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
		logrus.Debug("Creating an AWS Config:\n")
		if b.awsRegion != nil {
			logrus.Debugf("\tRegion: %s\n", *b.awsRegion)
		}
		if b.awsProfile != nil {
			logrus.Debugf("\tProfile: %s\n", *b.awsProfile)
		}
		if b.awsAccountID != nil {
			logrus.Debugf("\tAccountID: %s\n", *b.awsAccountID)
		}
		if b.awsRoleArn != nil {
			logrus.Debugf("\tRoleArn: %s\n", *b.awsRoleArn)
		}

		options := []func(*configv2.LoadOptions) error{
			configv2.WithRegion(*b.awsRegion),
			configv2.WithRetryer(func() aws.Retryer {
				// Unless specified, we run into ThrottlingException when repeating calls, when following logs or waiting on a condition.
				return retry.AddWithMaxBackoffDelay(retry.AddWithMaxAttempts(retry.NewStandard(), awsApiCallMaxRetries), awsApiCallBackoffDelay)
			}),
		}
		options = append(options, b.awsConfigLoadOptions...)

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

		if b.awsRoleArn != nil && len(*b.awsRoleArn) > 0 {
			stsClient := sts.NewFromConfig(conf)
			roleCreds := stscreds.NewAssumeRoleProvider(stsClient, *b.awsRoleArn)
			roleCfg := conf.Copy()
			roleCfg.Credentials = aws.NewCredentialsCache(roleCreds)
			conf = roleCfg
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
		logrus.Debugf("Creating an EKS client: region=%s. \n", b.awsConfig.Region)
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
		return nil, errors.Wrapf(err, "retrieving identity info")
	}
	logrus.Debugf("user identity confirmed: %s", userName)

	accountID, err := b.GetAccountID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve aws account id")
	}
	logrus.Debugf("AWS accunt ID confirmed: %s", accountID)

	_, err = b.GetComputeBackend(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to compute backend")
	}

	// other inferred or set fields
	if b.integrationSecret == nil {
		integrationSecret, integrationSecretArn, err := b.computeBackend.GetIntegrationSecret(ctx)
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
		EnvironmentContext: environmentContext,
		IntegrationSecret:  *b.integrationSecret,
	}

	return b, nil
}

func (b Backend) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	return b.awsConfig.Credentials.Retrieve(ctx)
}

func (b *Backend) GetComputeBackend(ctx context.Context) (compute.ComputeBackend, error) {
	var computeBackend compute.ComputeBackend
	var err error
	if b.environmentContext.TaskLaunchType == util.LaunchTypeK8S {
		computeBackend, err = NewK8SComputeBackend(ctx, b.environmentContext.K8S, b)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to connect to k8s backend")
		}
	} else if b.environmentContext.TaskLaunchType == util.LaunchTypeNull {
		computeBackend, err = NewNullComputeBackend(ctx, b)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to connect to null backend")
		}
	} else {
		computeBackend, err = NewECSComputeBackend(ctx, b.environmentContext.SecretID, b)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to connect to ecs backend")
		}
	}
	b.computeBackend = computeBackend
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

func (b *Backend) PrintLogs(ctx context.Context, stackName, serviceName, containerName string, opts ...util.PrintOption) error {
	return b.computeBackend.PrintLogs(ctx, stackName, serviceName, containerName, opts...)
}

func (b *Backend) RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error {
	return b.computeBackend.RunTask(ctx, taskDefArn, launchType)
}

func (b *Backend) Shell(ctx context.Context, stackName, serviceName, containerName string, shellCommand []string) error {
	return b.computeBackend.Shell(ctx, stackName, serviceName, containerName, shellCommand)
}

func (b *Backend) GetEvents(ctx context.Context, stackName string, services []string) error {
	return b.computeBackend.GetEvents(ctx, stackName, services)
}

func (b *Backend) Describe(ctx context.Context, stackName string, serviceName string) (compute.StackServiceDescription, error) {
	return b.computeBackend.Describe(ctx, stackName, serviceName)
}

func (b *Backend) ListEKSClusterIds(ctx context.Context) ([]string, error) {
	out, err := b.eksclient.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		return nil, err
	}
	return out.Clusters, nil
}

func (b *Backend) DisplayCloudWatchInsightsLink(ctx context.Context, logReference util.LogReference) error {
	queryId := uuid.NewUUID()
	cloudwatchLink, err := util.LogInsights2ConsoleLink(logReference,
		string(queryId))
	if err != nil {
		logrus.Errorf("To our dismay, we were unable to generate a link to query and visualize these logs")
	} else {
		if diagnostics.IsInteractiveContext(ctx) {
			proceed := false
			prompt := &survey.Confirm{Message: fmt.Sprintf("Would you like to query these logs in your browser? Please log into your AWS account (%s), then select Yes.", logReference.AWSAccountID)}
			err = survey.AskOne(prompt, &proceed)
			if err != nil || !proceed {
				return nil
			}
			logrus.Info("Opening Browser window to query cloudwatch insights.")
			err = browser.OpenURL(cloudwatchLink)
			if err != nil {
				return errors.Wrap(err, "To our dismay, we were unable open up a browser window to query cloudwatch insights.")
			}
			logrus.Info("Select the desired time frame, and click 'Run Query' to query the logs.")
			return nil
		}
		logrus.Info("****************************************************************************************")
		logrus.Infof("To query and visualize these logs, log into your AWS account (%s), navigate to the link below --", logReference.AWSAccountID)
		logrus.Info("(you will need to copy the entire link), and click 'Run Query' in AWS Console:")
		logrus.Info(cloudwatchLink)
		logrus.Info("****************************************************************************************")
		return nil
	}
	return nil
}
