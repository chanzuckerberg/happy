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
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/profiler"
	"github.com/chanzuckerberg/happy/pkg/util"
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

type RegistryDescriptor struct {
	RegistryId     string
	RepositoryName string
}

func (ab ArtifactBuilder) WithConfig(config *BuilderConfig) ArtifactBuilderIface {
	ab.config = config
	return ab
}

func (ab ArtifactBuilder) WithBackend(backend *backend.Backend) ArtifactBuilderIface {
	ab.backend = backend
	return ab
}

func (ab ArtifactBuilder) WithTags(tags []string) ArtifactBuilderIface {
	if len(tags) > 0 {
		ab.tags = tags
	}
	return ab
}

func (ab ArtifactBuilder) validate() error {
	if ab.config == nil {
		return errors.New("configuration was not provided")
	}
	if ab.backend == nil {
		return errors.New("backend was not provided")
	}
	return nil
}

func (ab ArtifactBuilder) CheckImageExists(ctx context.Context, tag string) (bool, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "CheckImageExists")
	err := ab.validate()
	if err != nil {
		return false, errors.Wrap(err, "artifact builder configuration is incomplete")
	}
	serviceRegistries := ab.backend.Conf().GetServiceRegistries()
	images, err := ab.config.GetBuildServicesImage(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get service image")
	}

	for serviceName := range images {
		registry, ok := serviceRegistries[serviceName]
		if !ok {
			continue
		}

		result, _, err := ab.getRegistryImages(ctx, registry, tag)
		if err != nil {
			return false, errors.Wrap(err, "error getting an image")
		}
		if result == nil || len(result.Images) == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (ab ArtifactBuilder) RetagImages(
	ctx context.Context,
	serviceRegistries map[string]*config.RegistryConfig,
	sourceTag string,
	destTags []string,
	images []string,
) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "RetagImages")
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

		log.Infof("retagging %s from '%s' to '%s'", serviceName, sourceTag, strings.Join(destTags, ","))

		result, descriptor, err := ab.getRegistryImages(ctx, registry, sourceTag)
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
				RepositoryName: aws.String(descriptor.RepositoryName),
				RegistryId:     aws.String(descriptor.RegistryId),
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

func (ab ArtifactBuilder) getRegistryImages(ctx context.Context, registry *config.RegistryConfig, tag string) (*ecr.BatchGetImageOutput, *RegistryDescriptor, error) {
	parts := strings.Split(registry.GetRepoUrl(), "/")
	if len(parts) < 2 {
		return nil, nil, errors.Errorf("invalid registry url format: %s", registry.GetRepoUrl())
	}
	registryId := parts[0]
	repositoryName := parts[1]

	if util.IsLocalstack() {
		registryId = "000000000000"
	} else {
		parts = strings.Split(registryId, ".")
		if len(parts) == 6 {
			// Real AWS registry ID
			registryId = parts[0]
		} else {
			return nil, nil, errors.Errorf("invalid registry format: %s", registryId)
		}
	}

	input := &ecr.BatchGetImageInput{
		ImageIds:           []ecrtypes.ImageIdentifier{{ImageTag: aws.String(tag)}},
		RepositoryName:     aws.String(repositoryName),
		AcceptedMediaTypes: []string{MediaTypeDocker1Manifest, MediaTypeDocker2Manifest, MediaTypeOCI1Manifest},
		RegistryId:         aws.String(registryId),
	}

	ecrClient := ab.backend.GetECRClient()
	result, err := ecrClient.BatchGetImage(ctx, input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error getting an image")
	}
	descriptor := RegistryDescriptor{RegistryId: registryId, RepositoryName: repositoryName}
	return result, &descriptor, nil
}

func (ab ArtifactBuilder) Build(ctx context.Context) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Build")

	_, err := ab.config.GetBuildServicesImage(ctx)
	if err != nil {
		return err
	}
	return ab.config.DockerComposeBuild()
}

func (ab ArtifactBuilder) RegistryLogin(ctx context.Context) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "RegistryLogin")
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

func (ab ArtifactBuilder) Push(ctx context.Context, tags []string) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Push")
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	serviceRegistries := ab.backend.Conf().GetServiceRegistries()
	servicesImage, err := ab.config.GetBuildServicesImage(ctx)
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
			// repo := registry.GetRepoUrl()
			// repo = strings.ReplaceAll(repo, ":4510/", "/")
			// image = strings.ReplaceAll(image, ":4510/", "/")
			dockerTagArgs := []string{"docker", "tag", fmt.Sprintf("%s:latest", image), fmt.Sprintf("%s:%s", registry.GetRepoUrl(), currentTag)}

			cmd := &exec.Cmd{
				Path:   docker,
				Args:   dockerTagArgs,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			log.Infof("executing: %s", cmd.String())

			if err := ab.config.executor.Run(cmd); err != nil {
				return errors.Wrap(err, "process failure")
			}

			// push image
			img := fmt.Sprintf("%s:%s", registry.GetRepoUrl(), currentTag)
			dockerPushArgs := []string{"docker", "push", img}

			cmd = &exec.Cmd{
				Path:   docker,
				Args:   dockerPushArgs,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			log.Infof("executing: %s", cmd.String())
			if err := ab.config.executor.Run(cmd); err != nil {
				return errors.Errorf("process failure: %v", err)
			}

		}
	}
	return nil
}

func (ab ArtifactBuilder) BuildAndPush(
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

	err = ab.Build(ctx)
	if err != nil {
		return err
	}

	return ab.Push(ctx, o.tags)
}
