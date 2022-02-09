package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type awsBackendOption func(*Backend)

// WithAWSRegion sets the AWS region for this Backend
func WithAWSRegion(region string) awsBackendOption {
	return func(ab *Backend) { ab.awsRegion = &region }
}

// WithAWSProfile sets the AWS profile to use for this Backend
func WithAWSProfile(profile string) awsBackendOption {
	return func(ab *Backend) { ab.awsProfile = profile }
}

// WithIntegrationSecret sets the IntegrationSecret for this Backend
func WithIntegrationSecret(integrationSecret *config.IntegrationSecret) awsBackendOption {
	return func(b *Backend) { b.integrationSecret = integrationSecret }
}

// WithSTSClients allows overriding the AWS STS Client
func WithSTSClient(client stsiface.STSAPI) awsBackendOption {
	return func(ab *Backend) { ab.stsclient = client }
}

// WithLogsClient allows overriding the AWS Logs Client
func WithLogsClient(client cloudwatchlogsiface.CloudWatchLogsAPI) awsBackendOption {
	return func(ab *Backend) { ab.logsclient = client }
}

// WithSSMClient allows overriding the AWS SSM Client
func WithSSMClient(client ssmiface.SSMAPI) awsBackendOption {
	return func(ab *Backend) { ab.ssmclient = client }
}

// WithECSClient allows overriding the AWS ECS Client
func WithECSClient(client ecsiface.ECSAPI) awsBackendOption {
	return func(ab *Backend) { ab.ecsclient = client }
}

// WithEC2Client allows overriding the AWS EC2 Client
func WithEC2Client(client ec2iface.EC2API) awsBackendOption {
	return func(ab *Backend) { ab.ec2client = client }
}

// WithECRClient allows overriding the AWS ECR Client
func WithECRClient(client ecriface.ECRAPI) awsBackendOption {
	return func(ab *Backend) { ab.ecrclient = client }
}

// WithSecretsClient allows overriding the AWS Secrets Client
func WithSecretsClient(client secretsmanageriface.SecretsManagerAPI) awsBackendOption {
	return func(ab *Backend) { ab.secretsclient = client }
}

// WithAWSSession allows configuring an AWS Session
func WithAWSSession(session *session.Session) awsBackendOption {
	return func(ab *Backend) { ab.awsSession = session }
}
