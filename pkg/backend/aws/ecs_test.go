package aws

import (
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNetworkConfig(t *testing.T) {
	r := require.New(t)
	backend := Backend{}
	sgs := []string{"sg-1", "sg-2"}
	subnets := []string{"subnet-1", "subnet-2"}

	backend.integrationSecret = &config.IntegrationSecret{
		ClusterArn:     "arn:cluster",
		PrivateSubnets: subnets,
		SecurityGroups: sgs,
		Services:       map[string]*config.RegistryConfig{},
	}
	networkConfig := backend.getNetworkConfig()
	r.NotNil(networkConfig)
	r.Equal(len(subnets), len(networkConfig.AwsvpcConfiguration.Subnets))
	r.Equal(len(sgs), len(networkConfig.AwsvpcConfiguration.SecurityGroups))

	for index, subnet := range subnets {
		r.Equal(subnet, networkConfig.AwsvpcConfiguration.Subnets[index])
	}
	for index, sg := range sgs {
		r.Equal(sg, networkConfig.AwsvpcConfiguration.SecurityGroups[index])
	}
}
