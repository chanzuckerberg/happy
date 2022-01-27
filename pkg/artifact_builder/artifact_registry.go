package artifact_builder

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

// import "strings"

type RegistryBackend interface {
	GetPwd(registryIds []string) (string, error)
	GetECRClient() *ecr.ECR
}
