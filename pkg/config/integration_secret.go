package config

import (
	"strings"

	"github.com/pkg/errors"
)

type RegistryConfig struct {
	Url string `json:"url"`
}

func (s *RegistryConfig) GetRepoUrl() string {
	return s.Url
}

func (s *RegistryConfig) GetRegistryUrl() string {
	return strings.Split(s.Url, "/")[0]
}

type TfeSecret struct {
	Url string `json:"url"`
	Org string `json:"org"`
}

type IntegrationSecret struct {
	ClusterArn     string                     `json:"cluster_arn"`
	PrivateSubnets []string                   `json:"private_subnets"`
	SecurityGroups []string                   `json:"security_groups"`
	Services       map[string]*RegistryConfig `json:"ecrs"`
	Tfe            *TfeSecret                 `json:"tfe"`
}

func (s *IntegrationSecret) GetClusterArn() string {
	return s.ClusterArn
}

func (s *IntegrationSecret) GetPrivateSubnets() []string {
	return s.PrivateSubnets
}

func (s *IntegrationSecret) GetSecurityGroups() []string {
	return s.SecurityGroups
}

func (s *IntegrationSecret) GetServiceUrl(serviceName string) (string, error) {
	svc, ok := s.Services[serviceName]
	if !ok {
		return "", errors.Errorf("can't find service %s", serviceName)
	}

	return svc.GetRepoUrl(), nil
}

func (s *IntegrationSecret) GetAllServicesUrl() map[string]*RegistryConfig {
	return s.Services
}

func (s *IntegrationSecret) GetTfeUrl() string {
	return s.Tfe.Url
}

func (s *IntegrationSecret) GetTfeOrg() string {
	return s.Tfe.Org
}

// // TODO(el): centralize the resolution of this into IntegrationSecret struct
// dockerRepo := os.Getenv("DOCKER_REPO")
// if len(dockerRepo) == 0 {
// 	serviceRegistries := secrets.GetAllServicesUrl()
// 	if err != nil {
// 		log.Errorf("Unable to retrieve registry information: %s\n", err.Error())
// 	}
// 	for _, registry := range serviceRegistries {
// 		dockerRepo = registry.Url
// 		parts := strings.Split(registry.GetRepoUrl(), "/")
// 		if len(parts) < 2 {
// 			continue
// 		}
// 		dockerRepo = parts[0] + "/"
// 		break
// 	}
// }
