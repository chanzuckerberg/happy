package artifact_builder

import (
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
)

var testFilePath = "../config/testdata/test_config.yaml"

func TestCheckTagExists(t *testing.T) {

	happyConfig, err := config.NewHappyConfig(testFilePath, "rdev")
	if err != nil {
		t.Error("Cannot load configuration")
		return
	}

	buildConfig := NewBuilderConfig("", "")
	artifactBuilder := NewArtifactBuilder(buildConfig, happyConfig)

	serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
	if err != nil {
		t.Error("Unable to get servce registries")
		return
	}

	err = artifactBuilder.CheckImageExists(serviceRegistries, "a")
	if err != nil {
		t.Error("Image check failed")
		return
	}
}
