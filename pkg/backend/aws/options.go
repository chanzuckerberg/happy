package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type awsBackendOption func(*awsBackend)

// WithAWSRegion sets the AWS region for this Backend
func WithAWSRegion(region string) awsBackendOption {
	return func(ab *awsBackend) { ab.awsRegion = &region }
}

// WithAWSProfile sets the AWS profile to use for this Backend
func WithAWSProfile(profile string) awsBackendOption {
	return func(ab *awsBackend) { ab.awsProfile = profile }
}

// WithSTSClients allows overriding the AWS STS Client
func WithSTSClient(client stsiface.STSAPI) awsBackendOption {
	return func(ab *awsBackend) { ab.stsclient = client }
}

// WithLogsClient allows overriding the AWS Logs Client
func WithLogsClient(client cloudwatchlogsiface.CloudWatchLogsAPI) awsBackendOption {
	return func(ab *awsBackend) { ab.logsclient = client }
}

// WithSSMClient allows overriding the AWS SSM Client
func WithSSMClient(client ssmiface.SSMAPI) awsBackendOption {
	return func(ab *awsBackend) { ab.ssmclient = client }
}

// WithECSClient allows overriding the AWS ECS Client
func WithECSClient(client ecsiface.ECSAPI) awsBackendOption {
	return func(ab *awsBackend) { ab.ecsclient = client }
}

// WithAWSSession allows configuring an AWS Session
func WithAWSSession(session *session.Session) awsBackendOption {
	return func(ab *awsBackend) { ab.awsSession = session }
}
