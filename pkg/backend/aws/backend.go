package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type awsBackend struct {
	// required inputs
	conf config.HappyConfig

	// optional inputs
	awsRegion  *string
	awsProfile string

	// aws settion: provided or inferred
	awsSession *session.Session

	// aws clients: provided or inferred
	stsclient  stsiface.STSAPI
	logsclient cloudwatchlogsiface.CloudWatchLogsAPI
	ssmclient  ssmiface.SSMAPI
	ecsclient  ecsiface.ECSAPI
}

// New returns a new AWS backend
func NewAWSBackend(conf config.HappyConfig, opts ...awsBackendOption) *awsBackend {
	// Set defaults
	b := &awsBackend{
		awsRegion: aws.String("us-west-2"),
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

	return b
}
