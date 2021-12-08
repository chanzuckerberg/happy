package artifact_builder

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

// import "strings"

type RegistryBackend interface {
	GetPwd(registryIds []string) (string, error)
	GetECRClient() *ecr.ECR
}

// type Registry struct {
// 	Url string `json:"url"`
// }

// func (s *Registry) GetRepoUrl() string {
// 	return s.Url
// }

// func (s *Registry) GetRegistryUrl() string {
// 	return strings.Split(s.Url, "/")[0]
// }
