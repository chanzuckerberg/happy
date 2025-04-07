package artifact_builder

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type DockerCommand string

const (
	DockerCommandConfig DockerCommand = "config"
	DockerCommandBuild  DockerCommand = "build"
)

func (bc *BuilderConfig) DockerComposeBuild() error {
	_, err := bc.invokeDockerCompose(DockerCommandBuild)
	return err
}

func (bc *BuilderConfig) DockerComposeConfig() (*ConfigData, error) {
	configDataBytes, err := bc.invokeDockerCompose(DockerCommandConfig)
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
func (bc *BuilderConfig) invokeDockerCompose(command DockerCommand) ([]byte, error) {
	composeArgs := []string{"docker", "compose", "--file", bc.composeFile}
	if len(bc.composeEnvFile) > 0 {
		composeArgs = append(composeArgs, "--env-file", bc.composeEnvFile)
	}

	// NOTE: by default this is the "*" (all) profile
	composeArgs = append(composeArgs, "--profile", bc.Profile.Get())

	envVars := bc.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)
	envVars = append(envVars, "DOCKER_BUILDKIT=1")

	docker, err := bc.Executor.LookPath("docker")
	if err != nil {
		return nil, errors.Wrap(err, "could not find docker compose in path")
	}

	cmd := &exec.Cmd{
		Path:   docker,
		Args:   append(composeArgs, string(command)),
		Env:    envVars,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
	logrus.Debugf("executing: %s", cmd.String())
	switch command {
	case DockerCommandConfig:
		output, err := bc.Executor.Output(cmd)
		return output, errors.Wrap(err, "unable to process docker compose output")
	default:
		cmd.Stdout = os.Stdout
		err = bc.Executor.Run(cmd)
		return []byte{}, errors.Wrap(err, "unable to process docker compose output")
	}
}
