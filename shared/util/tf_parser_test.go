package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTFParser(t *testing.T) {
	services, err := ParseServices("/Users/alokshin/GitHub/chanzuckerberg/k8s-test-app/.happy/terraform/envs/rdev")
	r := require.New(t)
	r.NoError(err)
	r.True(len(services) > 0)
}
