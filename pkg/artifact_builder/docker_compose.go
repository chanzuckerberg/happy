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

	dockerCompose, _ := exec.LookPath("docker-compose")

	cmd := &exec.Cmd{
		Path:   dockerCompose,
		Args:   append(composeArgs, command),
		Env:    envVars,
		Stderr: os.Stderr,
	}
	output, err := cmd.Output()

	if err != nil {
		return nil, errors.Wrap(err, "process failed:")
	}
	return output, nil
}
