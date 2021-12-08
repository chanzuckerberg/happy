package backend

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type UserIDBackend interface {
	GetUserName() (string, error)
}

type AwsSTS struct {
	session *session.Session
	awsConfig *aws.Config
	stsClient stsiface.STSAPI
}

var stsSessInst UserIDBackend
var creatStsOnce sync.Once

func GetAwsSts(config config.HappyConfigIface) UserIDBackend {
	awsProfile := config.AwsProfile()
	creatStsOnce.Do(func() {
		awsConfig := &aws.Config{
			Region: aws.String("us-west-2"),
			MaxRetries: aws.Int(2),
		}
		session := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: awsProfile,
			Config: *awsConfig,
			SharedConfigState: session.SharedConfigEnable,
		}))
		stsClient := sts.New(session)
		fmt.Println(stsClient)
		stsSessInst = &AwsSTS{
			session:         session,
			awsConfig:       awsConfig,
			stsClient: stsClient,
		}
	})
	return stsSessInst
}

func GetAwsStsWithClient(client stsiface.STSAPI) UserIDBackend {
	stsSessInst = &AwsSTS{
		stsClient: client,
	}
	return stsSessInst
}

func (s *AwsSTS) GetUserName() (string, error) {
	output, err := s.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	arnParts := strings.Split(*output.Arn, "/")
	return strings.Split(arnParts[len(arnParts)-1], "@")[0], nil
}
