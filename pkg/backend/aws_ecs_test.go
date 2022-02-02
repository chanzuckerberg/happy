package backend

import (
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/stretchr/testify/require"
)

var testFilePath = "../config/testdata/test_config.yaml"

func TestGetLogEvents(t *testing.T) {
	r := require.New(t)
	happyConfig, err := NewTestHappyConfig(t, testFilePath, "rdev")
	r.NoError(err)

	taskRunner := GetAwsEcs(happyConfig)
	ecsClient := taskRunner.GetECSClient()
	r.NotNil(ecsClient)
}

// generates a test happy config
// only use in tests
func NewTestHappyConfig(
	t *testing.T,
	testFilePath string,
	env string,
) (config.HappyConfig, error) {
	b := &config.Bootstrap{
		Env:             env,
		HappyConfigPath: testFilePath,
	}
	return config.NewHappyConfig(b)
}
