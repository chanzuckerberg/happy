package artifact_builder

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type DockerCommand string

const (
	DockerCommandConfig DockerCommand = "config"
	DockerCommandBuild  DockerCommand = "build"
)

func (bc *BuilderConfig) DockerComposeBuild(ctx context.Context) error {
	_, err := bc.invokeDockerCompose(ctx, DockerCommandBuild)
	return err
}

func (bc *BuilderConfig) DockerComposeConfig(ctx context.Context) (*ConfigData, error) {
	configDataBytes, err := bc.invokeDockerCompose(ctx, DockerCommandConfig)
	if err != nil {
		return nil, err
	}

	configData := &ConfigData{}
	err = yaml.Unmarshal(configDataBytes, configData)
	if err != nil {
		return nil, errors.Wrap(err, "could not yaml parse docker compose data")
	}
	return configData, nil
}

// 'docker-compose' was incorporated into 'docker' itself.
func (bc *BuilderConfig) invokeDockerCompose(ctx context.Context, command DockerCommand) ([]byte, error) {
	composeArgs := []string{"docker", "compose", "--file", bc.composeFile}
	if len(bc.composeEnvFile) > 0 {
		composeArgs = append(composeArgs, "--env-file", bc.composeEnvFile)
	}

	// NOTE: by default this is the "*" (all) profile
	composeArgs = append(composeArgs, "--profile", bc.Profile.Get())

	envVars := bc.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)
	envVars = append(envVars, "DOCKER_BUILDKIT=1")

	// LookPath is called by CommandContext
	cmd := exec.CommandContext(ctx, "docker", append(composeArgs, string(command))...)
	cmd.Env = envVars
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	logrus.Debugf("executing: %s", cmd.String())
	switch command {
	case DockerCommandConfig:
		output, err := cmd.Output()
		return output, errors.Wrap(err, "unable to process docker compose output")
	default:
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		return []byte{}, errors.Wrap(err, "unable to process docker compose output")
	}
}
