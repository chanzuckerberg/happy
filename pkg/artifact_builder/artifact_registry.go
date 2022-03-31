package artifact_builder

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type RegistryBackend interface {
	GetPwd(registryIds []string) (string, error)
	GetECRClient() *ecr.Client
}
