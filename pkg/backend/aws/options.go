package aws

import (
	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
func WithSTSClient(client *sts.Client) AWSBackendOption {
	return func(ab *Backend) { ab.stsclient = client }
}

func WithGetLogEventsAPIClient(client cwlv2.GetLogEventsAPIClient) AWSBackendOption {
	return func(ab *Backend) { ab.cwlGetLogEventsAPIClient = client }
}

// WithSSMClient allows overriding the AWS SSM Client
func WithSSMClient(client *ssm.Client) AWSBackendOption {
	return func(ab *Backend) { ab.ssmclient = client }
}

// WithECSClient allows overriding the AWS ECS Client
func WithECSClient(client *ecs.Client) AWSBackendOption {
	return func(ab *Backend) { ab.ecsclient = client }
}

// WithEC2Client allows overriding the AWS EC2 Client
func WithEC2Client(client *ec2.Client) AWSBackendOption {
	return func(ab *Backend) { ab.ec2client = client }
}

// WithECRClient allows overriding the AWS ECR Client
func WithECRClient(client *ecr.Client) AWSBackendOption {
	return func(ab *Backend) { ab.ecrclient = client }
}

// WithSecretsClient allows overriding the AWS Secrets Client
func WithSecretsClient(client *secretsmanager.Client) AWSBackendOption {
	return func(ab *Backend) { ab.secretsclient = client }
}

// WithAWSSession allows configuring an AWS Session
func WithAWSSession(config *aws2.Config) AWSBackendOption {
	return func(ab *Backend) { ab.awsSession = config }
}

// WithAWSAccountID allows configuring an AWS Account ID
func WithAWSAccountID(awsAccountID string) AWSBackendOption {
	return func(ab *Backend) { ab.awsAccountID = &awsAccountID }
}
