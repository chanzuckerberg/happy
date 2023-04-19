package artifact_builder

import (
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
)

type RegistryBackend interface {
	GetPwd(registryIds []string) (string, error)
	GetECRClient() interfaces.ECRAPI
}
