package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	compute "github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
)

type AWSBackendOption func(*Backend)

func WithNewAWSConfigOption(opt func(*configv2.LoadOptions) error) AWSBackendOption {
	return func(ab *Backend) { ab.awsConfigLoadOptions = append(ab.awsConfigLoadOptions, opt) }
}

// WithAWSRegion sets the AWS region for this Backend
func WithAWSRegion(region string) AWSBackendOption {
	return func(ab *Backend) { ab.awsRegion = &region }
}

// WithAWSProfile sets the AWS profile to use for this Backend
func WithAWSProfile(profile string) AWSBackendOption {
	return func(ab *Backend) { ab.awsProfile = &profile }
}

// WithAWSProfile sets the AWS profile to use for this Backend
func WithAWSRoleARN(awsRoleArn string) AWSBackendOption {
	return func(ab *Backend) { ab.awsRoleArn = &awsRoleArn }
}

// WithIntegrationSecret sets the IntegrationSecret for this Backend
func WithIntegrationSecret(integrationSecret *config.IntegrationSecret) AWSBackendOption {
	return func(b *Backend) { b.integrationSecret = integrationSecret }
}

// WithSTSClients allows overriding the AWS STS Client
func WithSTSClient(client interfaces.STSAPI) AWSBackendOption {
	return func(ab *Backend) { ab.stsclient = client }
}

// WithSTSClients allows overriding the AWS STS Client
func WithSTSPresignClient(client interfaces.STSPresignAPI) AWSBackendOption {
	return func(ab *Backend) { ab.stspresignclient = client }
}

func WithGetLogEventsAPIClient(client interfaces.GetLogEventsAPIClient) AWSBackendOption {
	return func(ab *Backend) { ab.cwlGetLogEventsAPIClient = client }
}

func WithFilterLogEventsAPIClient(client interfaces.FilterLogEventsAPIClient) AWSBackendOption {
	return func(ab *Backend) { ab.cwlFilterLogEventsAPIClient = client }
}

// WithSSMClient allows overriding the AWS SSM Client
func WithSSMClient(client interfaces.SSMAPI) AWSBackendOption {
	return func(ab *Backend) { ab.ssmclient = client }
}

// WithECSClient allows overriding the AWS ECS Client
func WithECSClient(client interfaces.ECSAPI) AWSBackendOption {
	return func(ab *Backend) { ab.ecsclient = client }
}

// WithEKSClient allows overriding the AWS EKS Client
func WithEKSClient(client interfaces.EKSAPI) AWSBackendOption {
	return func(ab *Backend) { ab.eksclient = client }
}

// WithEC2Client allows overriding the AWS EC2 Client
func WithEC2Client(client interfaces.EC2API) AWSBackendOption {
	return func(ab *Backend) { ab.ec2client = client }
}

// WithECRClient allows overriding the AWS ECR Client
func WithECRClient(client interfaces.ECRAPI) AWSBackendOption {
	return func(ab *Backend) { ab.ecrclient = client }
}

// WithSecretsClient allows overriding the AWS Secrets Client
func WithSecretsClient(client interfaces.SecretsManagerAPI) AWSBackendOption {
	return func(ab *Backend) { ab.secretsclient = client }
}

func WithDynamoDBClient(client interfaces.DynamoDB) AWSBackendOption {
	return func(ab *Backend) { ab.dynamodbclient = client }
}

// WithAWSConfig allows configuring an AWS Config
func WithAWSConfig(config *aws.Config) AWSBackendOption {
	return func(ab *Backend) { ab.awsConfig = config }
}

// WithAWSAccountID allows configuring an AWS Account ID
func WithAWSAccountID(awsAccountID string) AWSBackendOption {
	return func(ab *Backend) { ab.awsAccountID = &awsAccountID }
}

func WithTaskStoppedWaiter(waiter interfaces.ECSTaskStoppedWaiterAPI) AWSBackendOption {
	return func(ab *Backend) { ab.taskStoppedWaiter = waiter }
}

func WithK8SClientCreator(k8sClientCreator kube.K8sClientCreator) AWSBackendOption {
	return func(ab *Backend) { ab.k8sClientCreator = k8sClientCreator }
}

func WithComputeBackend(computeBackend compute.ComputeBackend) AWSBackendOption {
	return func(ab *Backend) { ab.computeBackend = computeBackend }
}

func WithExecutor(executor util.Executor) AWSBackendOption {
	return func(ab *Backend) { ab.executor = executor }
}
