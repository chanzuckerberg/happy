package artifact_builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ArtifactBuilder struct {
	backend *backend.Backend
	config  *BuilderConfig
}

func NewArtifactBuilder(builderConfig *BuilderConfig, backend *backend.Backend) *ArtifactBuilder {
	return &ArtifactBuilder{
		config:  builderConfig,
		backend: backend,
	}
}

func (s *ArtifactBuilder) CheckImageExists(serviceRegistries map[string]*config.RegistryConfig, tag string) (bool, error) {
	images, err := s.config.GetBuildServicesImage()
	if err != nil {
		return false, errors.Wrap(err, "failed to get service image")
	}

	for serviceName := range images {
		registry, ok := serviceRegistries[serviceName]
		if !ok {
			continue
		}

		parts := strings.Split(registry.GetRepoUrl(), "/")
		if len(parts) < 2 {
			return false, errors.Errorf("invalid registry url format: %s", registry.GetRepoUrl())
		}
		registryId := parts[0]
		repoUrl := parts[1]

		parts = strings.Split(registryId, ".")
		if len(parts) < 6 {
			return false, errors.Errorf("invalid registry id format: %s", registryId)
		}
		registryId = parts[0]

		ecrClient := s.backend.GetECRClient()

		input := &ecr.BatchGetImageInput{
			RegistryId: &registryId,
			ImageIds: []*ecr.ImageIdentifier{
				{
					ImageTag: aws.String(tag),
				},
			},
			RepositoryName: aws.String(repoUrl),
		}

		result, err := ecrClient.BatchGetImage(input)
		if err != nil {
			return false, errors.Wrap(err, "error getting an image")
		}
		if result == nil || len(result.Images) == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (s *ArtifactBuilder) RetagImages(serviceRegistries map[string]*config.RegistryConfig, servicesImage map[string]string, sourceTag string, destTags []string, images []string) error {
	ecrClient := s.backend.GetECRClient()

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
			fmt.Println("error Getting Image:", err)
			continue
		}

		if len(result.Images) == 0 {
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
				fmt.Println("error putting image", err)
				continue
			}
		}
	}

	return nil
}

func (s *ArtifactBuilder) Build() error {
	_, err := InvokeDockerCompose(*s.config, "build")
	if err != nil {
		return errors.Wrap(err, "build process failed:")
	}

	return nil
}

func (s *ArtifactBuilder) RegistryLogin(ctx context.Context) error {
	ecrAuthorizationToken, err := s.backend.ECRGetAuthorizationToken(ctx)
	if err != nil {
		return err
	}

	return ecrAuthorizationToken.DockerLogin(ctx)
}

func (s *ArtifactBuilder) Push(serviceRegistries map[string]*config.RegistryConfig, servicesImage map[string]string, tags []string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return errors.Wrap(err, "docker not in path")
	}
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
				return errors.Errorf("process failure: %v", err)
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
				return errors.Errorf("process failure: %v", err)
			}
		}
	}
	return nil
}
