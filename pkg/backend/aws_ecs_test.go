package backend

import (
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/stretchr/testify/require"
)

var testFilePath = "../config/testdata/test_config.yaml"

func TestGetLogEvents(t *testing.T) {
	r := require.New(t)
	happyConfig, err := config.NewHappyConfig(testFilePath, "rdev")
	r.NoError(err)

	taskRunner := GetAwsEcs(happyConfig)
	ecsClient := taskRunner.GetECSClient()
	r.NotNil(ecsClient)

}
