package backend

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

var ssmSessInst ParamStoreBackend
var creatSSMOnce sync.Once

type ParamStoreBackend interface {
	GetParameter(paramPath string) (*string, error)
	AddParams(name string, val string) error
}

type AwsSSMBackend struct {
	session   *session.Session
	config    *aws.Config
	ssmClient ssmiface.SSMAPI
}

func GetAwsBackend(config config.HappyConfigIface) ParamStoreBackend {
	awsProfile := config.AwsProfile()
	creatSSMOnce.Do(func() {
		config := &aws.Config{
			Region:     aws.String("us-west-2"),
			MaxRetries: aws.Int(2),
		}
		session := session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           awsProfile,
			Config:            *config,
			SharedConfigState: session.SharedConfigEnable,
		}))
		ssmClient := ssm.New(session)
		ssmSessInst = &AwsSSMBackend{
			session:   session,
			config:    config,
			ssmClient: ssmClient,
		}
	})
	return ssmSessInst
}

func GetAwsBackendWithClient(ssmClient ssmiface.SSMAPI) ParamStoreBackend {
	ssmSessInst = &AwsSSMBackend{
		ssmClient: ssmClient,
	}
	return ssmSessInst
}

func (s *AwsSSMBackend) GetParameter(paramPath string) (*string, error) {
	paramOutput, err := s.ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: &paramPath,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get params: %v", err)
	}

	return paramOutput.Parameter.Value, nil
}

func (s *AwsSSMBackend) AddParams(name string, val string) error {
	overwrite := true
	input := &ssm.PutParameterInput{
		Name:      &name,
		Overwrite: &overwrite,
		Value:     &val,
	}
	_, err := s.ssmClient.PutParameter(input)
	if err != nil {
		return err
	}
	return nil
}
