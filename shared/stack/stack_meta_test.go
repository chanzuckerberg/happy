package stack

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/shared/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../artifact_builder/testdata/test_config.yaml"
const testDockerComposePath = "../artifact_builder/testdata/docker-compose.yml"

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	r := require.New(t)
	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	config, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	stackMeta := &StackMeta{
		StackName: "test-stack",
		Env:       "rdev",
		Owner:     "test-owner",
	}

	// mock the backend
	backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext())
	r.NoError(err)

	username, err := backend.GetUserName(ctx)
	r.NoError(err)
	stackMeta.UpdateAll("test-tag", make(map[string]string), "", username, "/myapp", config, stackMeta.StackName, bootstrapConfig.Env)
	r.Equal(stackMeta.StackName, stackMeta.StackName)
	r.Equal(stackMeta.ImageTags, map[string]string{})
	r.Equal(stackMeta.Env, bootstrapConfig.Env)
	r.Equal(stackMeta.Owner, username)
	r.Equal(stackMeta.App, config.App())
	r.Equal(stackMeta.ImageTag, "test-tag")
	stackMeta.UpdateAll("test-tag", map[string]string{"foo": "bar"}, "", username, "/myapp", config, stackMeta.StackName, bootstrapConfig.Env)
	r.Equal(stackMeta.ImageTags, map[string]string{"foo": "bar"})
}
