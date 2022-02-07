package artifact_builder

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func InvokeDockerCompose(config BuilderConfig, command string) ([]byte, error) {
	composeArgs := []string{"docker-compose", "--file", config.composeFile}
	if len(config.envFile) > 0 {
		composeArgs = append(composeArgs, "--env-file", config.envFile)
	}

	envVars := config.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)

	dockerCompose, err := exec.LookPath("docker-compose")
	if err != nil {
		return nil, errors.Wrap(err, "could not find docker-compose in path")
	}

	cmd := &exec.Cmd{
		Path:   dockerCompose,
		Args:   append(composeArgs, command),
		Env:    envVars,
		Stderr: os.Stderr,
	}
	output, err := cmd.Output()
        return output, errors.Wrap(err, "process failed:")
}
