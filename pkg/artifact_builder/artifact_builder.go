package artifact_builder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ArtifactBuilder struct {
	config   *BuilderConfig
	registry RegistryBackend
}

func NewArtifactBuilder(builderConfig *BuilderConfig, happyConfig config.HappyConfigIface) *ArtifactBuilder {
	registry := GetECRBackend(happyConfig)
	return &ArtifactBuilder{
		config:   builderConfig,
		registry: registry,
	}
}

func (s *ArtifactBuilder) RetagImages(serviceRegistries map[string]*config.RegistryConfig, servicesImage map[string]string, sourceTag string, destTags []string, images []string) error {
	ecrClient := s.registry.GetECRClient()

	imageMap := make(map[string]bool)
	for _, image := range images {
		imageMap[image] = true
	}

	for serviceName, registry := range serviceRegistries {
		if _, ok := imageMap[serviceName]; !ok {
			if len(images) > 0 {
				continue
			}
		}

		repoUrl := strings.Split(registry.GetRepoUrl(), "/")[1]

		fmt.Printf("retagging %s from %s to %s", serviceName, sourceTag, destTags)

		input := &ecr.BatchGetImageInput{
			ImageIds: []*ecr.ImageIdentifier{
				{
					ImageTag: aws.String(sourceTag),
				},
			},
			RepositoryName: aws.String(repoUrl),
		}

		result, err := ecrClient.BatchGetImage(input)
		if err != nil {
			fmt.Println("Error Getting Image:", err)
			continue
		}

		manifest := result.Images[0].ImageManifest

		for _, tag := range destTags {
			input := &ecr.PutImageInput{
				ImageManifest:  manifest,
				ImageTag:       aws.String(tag),
				RepositoryName: aws.String(repoUrl),
			}

			_, err := ecrClient.PutImage(input)
			if err != nil {
				fmt.Println("Error putting image", err)
				continue
			}
		}
	}

	return nil
}

func (s *ArtifactBuilder) Build() error {
	composeArgs := []string{"docker-compose", "--file", s.config.composeFile}
	if s.config.env != "" {
		composeArgs = append(composeArgs, "--env", s.config.env)
	}

	envVars := s.config.GetBuildEnv()
	envVars = append(envVars, os.Environ()...)

	dockerCompose, _ := exec.LookPath("docker-compose")

	cmd := &exec.Cmd{
		Path:   dockerCompose,
		Args:   append(composeArgs, "build"),
		Env:    envVars,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Build process failed:")
	}

	return nil
}

func (s *ArtifactBuilder) RegistryLogin(serviceRegistries map[string]*config.RegistryConfig, images []string) error {
	registryIdSet := map[string]bool{}
	for _, registry := range serviceRegistries {
		regId := registry.GetRegistryUrl()
		if _, ok := registryIdSet[regId]; !ok {
			registryIdSet[regId] = true
		}
	}
	registryIds := []string{}
	for regId := range registryIdSet {
		registryIds = append(registryIds, regId)
	}
	registryPwd, err := s.registry.GetPwd(registryIds)
	if err != nil {
		return err
	}
	fmt.Println(registryIds)

	composeArgs := []string{"docker", "login", "--username", "AWS", "--password-stdin", registryIds[0]}

	docker, _ := exec.LookPath("docker")

	cmd := &exec.Cmd{
		Path:   docker,
		Args:   composeArgs,
		Stdin:  strings.NewReader(registryPwd),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Registry login failed:")
	}

	return nil
}

func (s *ArtifactBuilder) Push(serviceRegistries map[string]*config.RegistryConfig, servicesImage map[string]string, tags []string) error {

	docker, _ := exec.LookPath("docker")
	for serviceName, registry := range serviceRegistries {
		if _, ok := servicesImage[serviceName]; !ok {
			continue
		}

		image := servicesImage[serviceName]
		for _, currentTag := range tags {

			// re-tag image
			dockerTagArgs := []string{"docker", "tag", fmt.Sprintf("%s:latest", image), fmt.Sprintf("%s:%s", registry.GetRepoUrl(), currentTag)}
			log.WithField("args", dockerTagArgs).Debug("Running shell cmd")
			cmd := &exec.Cmd{
				Path:   docker,
				Args:   dockerTagArgs,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("Process failure: %v", err)
			}

			// push image
			dockerPushArgs := []string{"docker", "push", fmt.Sprintf("%s:%s", registry.GetRepoUrl(), currentTag)}
			log.WithField("args", dockerPushArgs).Debug("Running shell cmd")
			cmd = &exec.Cmd{
				Path:   docker,
				Args:   dockerPushArgs,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("Process failure: %v", err)
			}
		}
	}
	return nil
}
