package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type AWSBackendOption func(*Backend)

// WithAWSRegion sets the AWS region for this Backend
func WithAWSRegion(region string) AWSBackendOption {
	return func(ab *Backend) { ab.awsRegion = &region }
}

// WithAWSProfile sets the AWS profile to use for this Backend
func WithAWSProfile(profile string) AWSBackendOption {
	return func(ab *Backend) { ab.awsProfile = &profile }
}

// WithIntegrationSecret sets the IntegrationSecret for this Backend
func WithIntegrationSecret(integrationSecret *config.IntegrationSecret) AWSBackendOption {
	return func(b *Backend) { b.integrationSecret = integrationSecret }
}

// WithSTSClients allows overriding the AWS STS Client
func WithSTSClient(client interfaces.STSAPI) AWSBackendOption {
	return func(ab *Backend) { ab.stsclient = client }
}

func WithGetLogEventsAPIClient(client interfaces.GetLogEventsAPIClient) AWSBackendOption {
	return func(ab *Backend) { ab.cwlGetLogEventsAPIClient = client }
}

// WithSSMClient allows overriding the AWS SSM Client
func WithSSMClient(client interfaces.SSMAPI) AWSBackendOption {
	return func(ab *Backend) { ab.ssmclient = client }
}

// WithECSClient allows overriding the AWS ECS Client
func WithECSClient(client interfaces.ECSAPI) AWSBackendOption {
	return func(ab *Backend) { ab.ecsclient = client }
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

// WithAWSSession allows configuring an AWS Session
func WithAWSSession(config *aws.Config) AWSBackendOption {
	return func(ab *Backend) { ab.awsSession = config }
}

// WithAWSAccountID allows configuring an AWS Account ID
func WithAWSAccountID(awsAccountID string) AWSBackendOption {
	return func(ab *Backend) { ab.awsAccountID = &awsAccountID }
}
