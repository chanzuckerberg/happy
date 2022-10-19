package artifact_builder

import (
	"github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces"
)

type RegistryBackend interface {
	GetPwd(registryIds []string) (string, error)
	GetECRClient() interfaces.ECRAPI
}
