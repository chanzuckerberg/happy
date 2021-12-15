package backend

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/chanzuckerberg/happy-deploy/pkg/config"
)

type Cloudwatchlags struct {
	session      *session.Session
	awsConfig    *aws.Config
	awsLogClient cloudwatchlogsiface.CloudWatchLogsAPI
}

var cloudwatchlogsSessInst *Cloudwatchlags
var creatCloudwatchlogsOnce sync.Once

func GetAwsLogs(config config.HappyConfigIface) *Cloudwatchlags {
	awsProfile := config.AwsProfile()
	creatCloudwatchlogsOnce.Do(func() {
		awsConfig := &aws.Config{
			Region:     aws.String("us-west-2"),
			MaxRetries: aws.Int(2),
		}
		session := session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           awsProfile,
			Config:            *awsConfig,
			SharedConfigState: session.SharedConfigEnable,
		}))
		cloudwatchlogClient := cloudwatchlogs.New(session)
		cloudwatchlogsSessInst = &Cloudwatchlags{
			session:      session,
			awsLogClient: cloudwatchlogClient,
			awsConfig:    awsConfig,
		}
	})

	return cloudwatchlogsSessInst
}
