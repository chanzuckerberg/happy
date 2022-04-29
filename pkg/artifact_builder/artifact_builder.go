package artifact_builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/profiler"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ECR Supported media types
// https://docs.docker.com/registry/spec/manifest-v2-2/
// https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-manifest-formats.html
const (
	MediaTypeDocker1Manifest = "application/vnd.docker.distribution.manifest.v1+json"
	MediaTypeDocker2Manifest = "application/vnd.docker.distribution.manifest.v2+json"
	MediaTypeOCI1Manifest    = "application/vnd.oci.image.manifest.v1+json"
)

type ArtifactBuilder struct {
	backend  *backend.Backend
	config   *BuilderConfig
	Profiler *profiler.Profiler
	tags     []string
}

func NewArtifactBuilder() *ArtifactBuilder {
	return &ArtifactBuilder{
		config:   nil,
		backend:  nil,
		Profiler: profiler.NewProfiler(),
		tags:     []string{},
	}
}

func (ab *ArtifactBuilder) WithConfig(config *BuilderConfig) *ArtifactBuilder {
	ab.config = config
	return ab
}

func (ab *ArtifactBuilder) WithBackend(backend *backend.Backend) *ArtifactBuilder {
	ab.backend = backend
	return ab
}

func (ab *ArtifactBuilder) WithTags(tags []string) *ArtifactBuilder {
	if len(tags) > 0 {
		ab.tags = tags
	}
	return ab
}

func (ab *ArtifactBuilder) validate() error {
	if ab.config == nil {
		return errors.New("configuration was not provided")
	}
	if ab.backend == nil {
		return errors.New("backend was not provided")
	}
	return nil
}

func (ab *ArtifactBuilder) CheckImageExists(ctx context.Context, tag string) (bool, error) {
	defer ab.Profiler.AddRuntime(time.Now(), "CheckImageExists")
	err := ab.validate()
	if err != nil {
		return false, errors.Wrap(err, "artifact builder configuration is incomplete")
	}
	serviceRegistries := ab.backend.Conf().GetServiceRegistries()
	images, err := ab.config.GetBuildServicesImage()
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

		ecrClient := ab.backend.GetECRClient()

		input := &ecr.BatchGetImageInput{
			ImageIds:           []ecrtypes.ImageIdentifier{{ImageTag: aws.String(tag)}},
			RepositoryName:     aws.String(repoUrl),
			AcceptedMediaTypes: []string{MediaTypeDocker1Manifest, MediaTypeDocker2Manifest, MediaTypeOCI1Manifest},
			RegistryId:         &registryId,
		}

		result, err := ecrClient.BatchGetImage(ctx, input)
		if err != nil {
			return false, errors.Wrapf(err, "error getting an image (%s:%s)", repoUrl, tag)
		}
		if result == nil || len(result.Images) == 0 {
			return false, errors.Errorf("image (%s:%s) not found", repoUrl, tag)
		}
	}

	return true, nil
}

func (ab *ArtifactBuilder) RetagImages(
	ctx context.Context,
	serviceRegistries map[string]*config.RegistryConfig,
	sourceTag string,
	destTags []string,
	images []string,
) error {
	defer ab.Profiler.AddRuntime(time.Now(), "RetagImages")
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	ecrClient := ab.backend.GetECRClient()

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

		log.Infof("retagging %s from '%s' to '%s'", serviceName, sourceTag, strings.Join(destTags, ","))

		input := &ecr.BatchGetImageInput{
			ImageIds:           []ecrtypes.ImageIdentifier{{ImageTag: aws.String(sourceTag)}},
			RepositoryName:     aws.String(repoUrl),
			AcceptedMediaTypes: []string{},
			RegistryId:         new(string),
		}

		result, err := ecrClient.BatchGetImage(ctx, input)
		if err != nil {
			log.Errorf("error getting Image: %s", err.Error())
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

			_, err := ecrClient.PutImage(ctx, input)
			if err != nil {
				log.Error("error putting image", err)
				continue
			}
		}
	}

	return nil
}

func (ab *ArtifactBuilder) Build() error {
	defer ab.Profiler.AddRuntime(time.Now(), "Build")
	return ab.config.DockerComposeBuild()
}

func (ab *ArtifactBuilder) RegistryLogin(ctx context.Context) error {
	defer ab.Profiler.AddRuntime(time.Now(), "RegistryLogin")
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	ecrAuthorizationToken, err := ab.backend.ECRGetAuthorizationToken(ctx)
	if err != nil {
		return err
	}

	args := []string{"login", "--username", ecrAuthorizationToken.Username, "--password", ecrAuthorizationToken.Password, ecrAuthorizationToken.ProxyEndpoint}

	docker, err := ab.config.executor.LookPath("docker")
	if err != nil {
		return errors.Wrap(err, "could not find docker in path")
	}
	cmd := exec.CommandContext(ctx, docker, args...)
	err = ab.config.executor.Run(cmd)
	return errors.Wrap(err, "registry login failed")
}

func (ab *ArtifactBuilder) Push(tags []string) error {
	defer ab.Profiler.AddRuntime(time.Now(), "Push")
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	serviceRegistries := ab.backend.Conf().GetServiceRegistries()
	servicesImage, err := ab.config.GetBuildServicesImage()
	if err != nil {
		return err
	}

	docker, err := ab.config.executor.LookPath("docker")
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
			if err := ab.config.executor.Run(cmd); err != nil {
				return errors.Wrap(err, "process failure")
			}

			// push image
			img := fmt.Sprintf("%s:%s", registry.GetRepoUrl(), currentTag)
			dockerPushArgs := []string{"docker", "push", img}
			log.WithField("args", dockerPushArgs).Debug("Running shell cmd")
			cmd = &exec.Cmd{
				Path:   docker,
				Args:   dockerPushArgs,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			if err := ab.config.executor.Run(cmd); err != nil {
				return errors.Errorf("process failure: %v", err)
			}
			log.WithField("args", dockerTagArgs).Info("Tagged the image")
		}
	}
	return nil
}

func (ab *ArtifactBuilder) BuildAndPush(
	ctx context.Context,
	opts ...ArtifactBuilderBuildOption,
) error {
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	// calculate defaults
	defaultTag, err := ab.backend.GenerateTag(ctx)
	if err != nil {
		return err
	}
	tags := []string{defaultTag}
	if len(ab.tags) > 0 {
		tags = append(tags, ab.tags...)
	}

	// Get all the options first
	o := &artifactBuilderBuildOptions{
		tags: tags,
	}
	for _, opt := range opts {
		opt(o)
	}

	// Run logic
	err = ab.RegistryLogin(ctx)
	if err != nil {
		return err
	}

	err = ab.Build()
	if err != nil {
		return err
	}

	return ab.Push(o.tags)
}
