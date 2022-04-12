package artifact_builder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type DockerCommand string

const (
	DockerCommandConfig         DockerCommand = "config"
	DockerCommandBuild          DockerCommand = "build"
	DockerDefaultPlatformEnvVar string        = "DOCKER_DEFAULT_PLATFORM"
)

func (bc *BuilderConfig) DockerComposeBuild(targetPlatformDefinedInDockerCompose bool) error {
	_, err := bc.invokeDockerCompose(DockerCommandBuild, targetPlatformDefinedInDockerCompose)
	return err
}

func (bc *BuilderConfig) DockerComposeConfig() (*ConfigData, error) {
	configDataBytes, err := bc.invokeDockerCompose(DockerCommandConfig, false)
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
func (bc *BuilderConfig) invokeDockerCompose(command DockerCommand, targetPlatformDefinedInDockerCompose bool) ([]byte, error) {
	composeArgs := []string{"docker", "compose", "--file", bc.composeFile}
	if len(bc.composeEnvFile) > 0 {
		composeArgs = append(composeArgs, "--env-file", bc.composeEnvFile)
	}

	// NOTE: by default this is the "*" (all) profile
	composeArgs = append(composeArgs, "--profile", bc.profile.Get())

	envVars := bc.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)

	// Specifying platform in the docker compose conflicts with the DOCKER_DEFAULT_PLATFORM env var:
	// multiple platforms feature is currently not supported for docker driver. Please switch to a different driver (eg. "docker buildx create --use")
	envVars = filterOutTargetPlatformEnv(envVars)
	if !targetPlatformDefinedInDockerCompose {
		if bc.targetContainerPlatform != util.GetUserContainerPlatform() {
			envVars = append(envVars, fmt.Sprintf("%s=%s", DockerDefaultPlatformEnvVar, bc.targetContainerPlatform))
		}
	}

	docker, err := bc.GetExecutor().LookPath("docker")
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
	log.Infof("executing: %s", cmd.String())

	switch command {
	case DockerCommandConfig:
		output, err := bc.GetExecutor().Output(cmd)
		return output, errors.Wrap(err, "unable to process docker compose output")
	default:
		cmd.Stdout = os.Stdout
		err = bc.GetExecutor().Run(cmd)
		return []byte{}, errors.Wrap(err, "unable to process docker compose output")
	}
}

func filterOutTargetPlatformEnv(envVars []string) []string {
	filteredEnvVars := []string{}
	for _, envVar := range envVars {
		if !strings.Contains(envVar, DockerDefaultPlatformEnvVar) {
			filteredEnvVars = append(filteredEnvVars, envVar)
		}
	}
	return filteredEnvVars
}
