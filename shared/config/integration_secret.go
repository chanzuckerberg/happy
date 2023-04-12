package config

import (
	"github.com/pkg/errors"
)

type RegistryConfig struct {
	URL string `json:"url"`
}

type TfeSecret struct {
	Url string `json:"url"`
	Org string `json:"org"`
}

type IntegrationSecret struct {
	ClusterArn     string   `json:"cluster_arn"`
	PrivateSubnets []string `json:"private_subnets"`
	SecurityGroups []string `json:"security_groups"`
	// HACK HACK:
	//                 - We alias `ecrs` to Services. Misdirection, be careful..
	//                 - We only push images where the dockercompose.<service_name> match the ecr <registry name>. Otherwise we skip.
	Services            map[string]*RegistryConfig `json:"ecrs"`
	Tfe                 *TfeSecret                 `json:"tfe"`
	DynamoLocktableName string                     `json:"dynamo_locktable_name"`
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

	return svc.URL, nil
}

func (s *IntegrationSecret) GetServiceRegistries() map[string]*RegistryConfig {
	return s.Services
}

func (s *IntegrationSecret) GetTfeUrl() string {
	return s.Tfe.Url
}

func (s *IntegrationSecret) GetTfeOrg() string {
	return s.Tfe.Org
}

func (s *IntegrationSecret) GetDynamoLocktableName() string {
	return s.DynamoLocktableName
}
