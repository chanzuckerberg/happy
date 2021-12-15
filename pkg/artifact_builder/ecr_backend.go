package artifact_builder

import (
	"strings"
	"sync"

	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/chanzuckerberg/happy-deploy/pkg/config"
)

type EcrBackend struct {
	session   *session.Session
	config    *aws.Config
	ecrClient *ecr.ECR
}

var ecrSessInst *EcrBackend
var creatECROnce sync.Once

func GetECRBackend(config config.HappyConfigIface) RegistryBackend {
	awsProfile := config.AwsProfile()
	creatECROnce.Do(func() {
		config := &aws.Config{
			Region:     aws.String("us-west-2"),
			MaxRetries: aws.Int(2),
		}
		session := session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           awsProfile,
			Config:            *config,
			SharedConfigState: session.SharedConfigEnable,
		}))
		ecrClient := ecr.New(session)
		ecrSessInst = &EcrBackend{
			session:   session,
			config:    config,
			ecrClient: ecrClient,
		}
	})
	return ecrSessInst
}

func (s *EcrBackend) GetPwd(registryIds []string) (string, error) {
	registryIdsPtr := make([]*string, len(registryIds))

	for i, regId := range registryIds {
		registryIdsPtr[i] = &regId
	}
	tokens, err := s.ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}
	tokenDec, _ := b64.StdEncoding.DecodeString(*tokens.AuthorizationData[0].AuthorizationToken)
	return strings.Split(string(tokenDec), ":")[1], nil
}

func (s *EcrBackend) GetECRClient() *ecr.ECR {
	return s.ecrClient
}
