package backend

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type Cloudwatchlogs struct {
	session      *session.Session
	awsConfig    *aws.Config
	awsLogClient cloudwatchlogsiface.CloudWatchLogsAPI
}

var cloudwatchlogsSessInst *Cloudwatchlogs
var creatCloudwatchlogsOnce sync.Once

func GetAwsLogs(config config.HappyConfig) *Cloudwatchlogs {
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
		cloudwatchlogsSessInst = &Cloudwatchlogs{
			session:      session,
			awsLogClient: cloudwatchlogClient,
			awsConfig:    awsConfig,
		}
	})

	return cloudwatchlogsSessInst
}
