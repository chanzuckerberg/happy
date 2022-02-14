package artifact_builder

import (
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// 'docker-compose' was incorporated into 'docker' itself.
func InvokeDockerCompose(config BuilderConfig, command string) ([]byte, error) {
	composeArgs := []string{"docker", "compose", "--file", config.composeFile}
	if len(config.envFile) > 0 {
		composeArgs = append(composeArgs, "--env-file", config.envFile)
	}

	envVars := config.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)
	envVars = append(envVars, "DOCKER_BUILDKIT=0")

	docker, err := exec.LookPath("docker")
	if err != nil {
		return nil, errors.Wrap(err, "could not find docker-compose in path")
	}

	cmd := &exec.Cmd{
		Path:   docker,
		Args:   append(composeArgs, command),
		Env:    envVars,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
	log.Printf("Executing: %s\n", cmd.String())
	if command == "config" {
		output, err := config.GetExecutor().Output(cmd)
		return output, errors.Wrap(err, "unable to process docker compose output")
	} else {
		cmd.Stdout = os.Stdout
		err = config.GetExecutor().Run(cmd)
		return []byte{}, errors.Wrap(err, "unable to process docker compose output")
	}
}
