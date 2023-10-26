package artifact_builder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/profiler"
	stackservice "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
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
	backend     *backend.Backend
	config      *BuilderConfig
	happyConfig *config.HappyConfig
	Profiler    *profiler.Profiler
	tags        []string
}

type RegistryDescriptor struct {
	RegistryId     string
	RepositoryName string
}

func (ab ArtifactBuilder) GetTags() []string {
	return ab.tags
}

func (ab *ArtifactBuilder) WithConfig(config *BuilderConfig) ArtifactBuilderIface {
	ab.config = config
	return ab
}

func (ab *ArtifactBuilder) WithBackend(backend *backend.Backend) ArtifactBuilderIface {
	ab.backend = backend
	return ab
}

func (ab *ArtifactBuilder) WithHappyConfig(happyConfig *config.HappyConfig) ArtifactBuilderIface {
	ab.happyConfig = happyConfig
	return ab
}

func (ab *ArtifactBuilder) WithTags(tags []string) ArtifactBuilderIface {
	t := []string{}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		t = append(t, tag)
	}
	ab.tags = t
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
	// backward compatible way of overriding the new ECR locations
	// if they exist. The new ECRs will look like <stackname>/<env>/<servicename>
	// if users haven't switched to the latest stack TF module, this will return nothing
	stackECRS, err := ab.GetECRsForServices(ctx)
	if err != nil {
		log.Debugf("unable to get ECRs for services: %s", err)
	}
	if len(stackECRS) > 0 && err == nil {
		serviceRegistries = stackECRS
	}
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

		return true, nil
	}

	return false, nil
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
	validImageMap := make(map[string]bool)
	for _, image := range ab.happyConfig.GetServices() {
		validImageMap[image] = true
	}

	imageMap := make(map[string]bool)
	for _, image := range images {
		imageMap[image] = true
	}

	if len(images) == 0 {
		imageMap = validImageMap
	}

	for serviceName, registry := range serviceRegistries {
		if _, ok := validImageMap[serviceName]; !ok {
			continue
		}
		if _, ok := imageMap[serviceName]; !ok {
			continue
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
	parts := strings.SplitN(registry.URL, "/", 2)
	if len(parts) < 2 {
		return nil, nil, errors.Errorf("invalid registry url format: %s", registry.URL)
	}
	registryId := parts[0]
	repositoryName := parts[1]

	if util.IsLocalstackMode() {
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

	docker, err := ab.config.Executor.LookPath("docker")
	if err != nil {
		return errors.Wrap(err, "could not find docker in path")
	}
	cmd := exec.CommandContext(ctx, docker, args...)
	err = ab.config.Executor.Run(cmd)
	return errors.Wrap(err, "registry login failed")
}

func (ab *ArtifactBuilder) getStacksECRSFromTFE(ctx context.Context, tfeWorkspace workspace_repo.Workspace, stackName string) (map[string]*config.RegistryConfig, error) {
	outs, err := tfeWorkspace.GetOutputs(ctx)
	if err != nil {
		if ab.happyConfig.GetData().FeatureFlags.EnableECRAutoCreation {
			log.Errorf("unable to get state outputs from stack workspace %s, cannot determine which ECR repo to push to", stackName)
		}
	}

	serviceECRs, ok := outs["service_ecrs"]
	if !ok {
		if ab.happyConfig.GetData().FeatureFlags.EnableECRAutoCreation {
			log.Errorf("unable to get service_ecrs from stack outputs %s, cannot determine which ECR repo to push to", stackName)
		}
		return nil, nil
	}

	sr := map[string]string{}
	err = json.Unmarshal([]byte(serviceECRs), &sr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal TFE workspace output")
	}

	serviceRegistries := map[string]*config.RegistryConfig{}
	for k, v := range sr {
		serviceRegistries[k] = &config.RegistryConfig{
			URL: v,
		}
	}
	return serviceRegistries, nil
}

func (ab ArtifactBuilder) GetECRsForServices(ctx context.Context) (map[string]*config.RegistryConfig, error) {
	repo := workspace_repo.NewWorkspaceRepo(ab.backend.Conf().GetTfeUrl(), ab.backend.Conf().GetTfeOrg())
	stackService := stackservice.NewStackService(ab.happyConfig.GetEnv(), ab.happyConfig.App()).WithBackend(ab.backend).WithWorkspaceRepo(repo)
	tfeWorkspace, err := stackService.GetStackWorkspace(ctx, ab.config.StackName)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get workspace for stack %s", ab.config.StackName)
	}

	return ab.getStacksECRSFromTFE(ctx, tfeWorkspace, ab.config.StackName)
}

func (ab *ArtifactBuilder) Pull(ctx context.Context, stackName, tag string) (map[string]string, error) {
	if tag == "" {
		return nil, errors.New("when pulling an image, the tag is required since we don't support a default tag")
	}
	err := ab.RegistryLogin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login to registry")
	}

	repo := workspace_repo.NewWorkspaceRepo(ab.backend.Conf().GetTfeUrl(), ab.backend.Conf().GetTfeOrg())
	tfeWorkspace, err := repo.GetWorkspace(ctx, fmt.Sprintf("%s-%s", ab.config.env, stackName))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get workspace for stack %s-%s", ab.config.env, stackName)
	}

	serviceRegistries, err := ab.getStacksECRSFromTFE(ctx, tfeWorkspace, stackName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get service registries from TFE; this feature requires autocreated ECRs")
	}

	servicesImage := map[string]string{}
	for service, registry := range serviceRegistries {
		dest := fmt.Sprintf("%s:%s", registry.URL, tag)
		cmd := exec.CommandContext(ctx, "docker", "pull", dest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := ab.config.Executor.Run(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "error running docker pull")
		}
		servicesImage[service] = dest
	}
	return servicesImage, nil
}

func (ab ArtifactBuilder) push(ctx context.Context, tags []string, servicesImage map[string]string) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Push")
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	serviceRegistries := ab.backend.Conf().GetServiceRegistries()

	// backward compatible way of overriding the new ECR locations
	// if they exist. The new ECRs will look like <stackname>-<servicename>
	stackECRS, err := ab.GetECRsForServices(ctx)
	if err != nil {
		log.Debugf("unable to get ECRs for services: %s", err)
	}
	if len(stackECRS) > 0 && err == nil {
		serviceRegistries = stackECRS
	}

	for serviceName, registry := range serviceRegistries {
		if _, ok := servicesImage[serviceName]; !ok {
			continue
		}

		image := servicesImage[serviceName]
		for _, currentTag := range tags {
			// re-tag image
			cmd := exec.CommandContext(ctx, "docker", "tag", image, fmt.Sprintf("%s:%s", registry.URL, currentTag))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := ab.config.Executor.Run(cmd)
			if err != nil {
				return errors.Wrap(err, "error tagging docker image")
			}

			// push image
			cmd = exec.CommandContext(ctx, "docker", "push", fmt.Sprintf("%s:%s", registry.URL, currentTag))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = ab.config.Executor.Run(cmd)
			if err != nil {
				return errors.Wrap(err, "error pushing docker image")
			}
		}
	}

	ab.cveScan(ctx, serviceRegistries, servicesImage, tags)

	return nil
}

type repositoryConfig struct {
	scanOnPush  bool
	descriptor  *RegistryDescriptor
	images      []ecrtypes.Image
	serviceName string
	tag         string
}

func (ab ArtifactBuilder) cveScan(ctx context.Context, serviceRegistries map[string]*config.RegistryConfig, servicesImage map[string]string, tags []string) {
	ecrClient := ab.backend.GetECRClient()
	scanConfiguration := []repositoryConfig{}
	repoScanConfiguration := map[string]bool{}

	globalScanEnabled := false

	out, err := ecrClient.GetRegistryScanningConfiguration(ctx, &ecr.GetRegistryScanningConfigurationInput{})
	if err != nil {
		log.Errorf("Unable to get registry scan configuration: %s", err.Error())
	} else if len(out.ScanningConfiguration.Rules) > 0 {
		for _, rule := range out.ScanningConfiguration.Rules {
			if rule.ScanFrequency == ecrtypes.ScanFrequencyScanOnPush {
				for _, filter := range rule.RepositoryFilters {
					if filter.FilterType == ecrtypes.ScanningRepositoryFilterTypeWildcard && filter.Filter != nil && *filter.Filter == "*" {
						globalScanEnabled = true
						break
					}
				}
			}
			if globalScanEnabled {
				break
			}
		}
	}

	scanningEnabled := globalScanEnabled

	if globalScanEnabled {
		log.Info("Your ECR registry has global scanning on push enabled.")
	}

	for serviceName, registry := range serviceRegistries {
		if _, ok := servicesImage[serviceName]; !ok {
			continue
		}
		for _, currentTag := range tags {
			result, descriptor, err := ab.getRegistryImages(ctx, registry, currentTag)
			if err != nil {
				log.Errorf("error getting Image: %s", err.Error())
				continue
			}

			repoScanEnabled, ok := repoScanConfiguration[descriptor.RepositoryName]

			if !ok {
				repoScanEnabled = globalScanEnabled

				if !repoScanEnabled {
					out, err := ecrClient.BatchGetRepositoryScanningConfiguration(ctx, &ecr.BatchGetRepositoryScanningConfigurationInput{
						RepositoryNames: []string{descriptor.RepositoryName},
					})

					if err != nil {
						log.Errorf("error getting repository scanning configuration: %s", err.Error())
					} else if len(out.ScanningConfigurations) > 0 {
						repoScanEnabled = out.ScanningConfigurations[0].ScanOnPush
					}
				}

				repoScanConfiguration[descriptor.RepositoryName] = repoScanEnabled
			}

			scanningEnabled = scanningEnabled || repoScanEnabled

			repositoryConfig := repositoryConfig{
				images:      result.Images,
				descriptor:  descriptor,
				tag:         currentTag,
				serviceName: serviceName,
				scanOnPush:  repoScanEnabled,
			}

			scanConfiguration = append(scanConfiguration, repositoryConfig)
		}
	}

	if scanningEnabled {
		log.Info("Scanning images for vulnerabilities...")
	} else {
		log.Warn("Scanning is not enabled on any of your services. Please consider adding 'scan_on_push = \"true\"' to your service definitions.")
		return
	}

	vulnerabilityCount := 0

	for _, repoConfig := range scanConfiguration {
		descriptor := repoConfig.descriptor
		for _, image := range repoConfig.images {
			imageRef := fmt.Sprintf("%s/%s:%s", *image.RegistryId, *image.RepositoryName, *image.ImageId.ImageTag)
			log.Infof("Waiting for service '%s' image %s ECR scan to complete\n", repoConfig.serviceName, imageRef)

			waiter := ecr.NewImageScanCompleteWaiter(ecrClient)
			err := waiter.Wait(ctx, &ecr.DescribeImageScanFindingsInput{
				RegistryId:     &descriptor.RegistryId,
				RepositoryName: &descriptor.RepositoryName,
				ImageId: &types.ImageIdentifier{
					ImageDigest: image.ImageId.ImageDigest,
					ImageTag:    image.ImageId.ImageTag,
				},
			}, 120*time.Second, func(opts *ecr.ImageScanCompleteWaiterOptions) {
				opts.LogWaitAttempts = true
			})

			if err != nil {
				log.Errorf("error waiting for image scan: %s", err.Error())
				continue
			}

			config := ecr.DescribeImageScanFindingsInput{
				RepositoryName: &descriptor.RepositoryName,
				RegistryId:     &descriptor.RegistryId,
				ImageId: &ecrtypes.ImageIdentifier{
					ImageDigest: image.ImageId.ImageDigest,
					ImageTag:    image.ImageId.ImageTag,
				},
				MaxResults: aws.Int32(1000),
			}

			paginator := ecr.NewDescribeImageScanFindingsPaginator(ecrClient, &config)

			for paginator.HasMorePages() {
				res, err := paginator.NextPage(ctx)
				if err != nil {
					log.Errorf("error getting image scan findings: %s", err.Error())
					continue
				}
				if res.ImageScanFindings == nil {
					continue
				}

				for _, finding := range res.ImageScanFindings.Findings {
					if finding.Severity == ecrtypes.FindingSeverityHigh || finding.Severity == ecrtypes.FindingSeverityCritical {
						vulnerabilityCount++
						log.Warnf("[%s] Vulnerability found in %s: %s %s", finding.Severity, imageRef, *finding.Name, *finding.Uri)
					}
				}
			}
		}
	}

	if vulnerabilityCount == 0 {
		log.Info("No high or critical vulnerabilities were found.")
		return
	}
	log.Warnf("Found %d high or critical vulnerabilities.", vulnerabilityCount)
}

// Push takes the source images from the docker compose file and uses "latest" tag
// of the built images on your local machine to push to the repository
func (ab ArtifactBuilder) Push(ctx context.Context, tags []string) error {
	servicesImage, err := ab.config.GetBuildServicesImage(ctx)
	if err != nil {
		return err
	}
	for k, v := range servicesImage {
		servicesImage[k] = fmt.Sprintf("%s:%s", v, "latest")
	}
	return ab.push(ctx, tags, servicesImage)
}

// PushFrom allows the caller to specify where the images are coming from and also what tags
// to pull from.
func (ab ArtifactBuilder) PushFromWithTag(ctx context.Context, servicesImage map[string]string, tag string) error {
	err = ab.RegistryLogin(ctx)
	if err != nil {
		return err
	}
	
	return ab.push(ctx, []string{tag}, servicesImage)
}

func (ab ArtifactBuilder) BuildAndPush(ctx context.Context) error {
	err := ab.validate()
	if err != nil {
		return errors.Wrap(err, "artifact builder configuration is incomplete")
	}

	err = ab.RegistryLogin(ctx)
	if err != nil {
		return err
	}

	err = ab.Build(ctx)
	if err != nil {
		return err
	}

	return ab.Push(ctx, ab.tags)
}

func (ab ArtifactBuilder) GetServices(ctx context.Context) (map[string]ServiceConfig, error) {
	config, err := ab.config.GetConfigData(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config data")
	}
	return config.Services, nil
}

func (ab ArtifactBuilder) GetAllServices(ctx context.Context) (map[string]ServiceConfig, error) {
	bc := ab.config.Clone()
	bc.configData = nil
	bc.Profile = nil
	config, err := bc.GetConfigData(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config data")
	}
	return config.Services, nil
}
