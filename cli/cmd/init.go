package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	ab "github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	waitoptions "github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type HappyClient struct {
	HappyConfig     *config.HappyConfig
	StackService    *stackservice.StackService
	ArtifactBuilder ab.ArtifactBuilderIface
	StackTags       map[string]string
	AWSBackend      *backend.Backend
}

func makeHappyClient(cmd *cobra.Command, sliceName, stackName string, tags []string, createTag, dryRun bool) (*HappyClient, error) {
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	ctx := cmd.Context()
	awsBackend, err := backend.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
	if err != nil {
		return nil, err
	}
	builderConfig := ab.
		NewBuilderConfig().
		WithBootstrap(bootstrapConfig).
		WithHappyConfig(happyConfig)

	builderConfig.DryRun = dryRun
	builderConfig.StackName = stackName
	ab, stackTags, err := configureArtifactBuilder(ctx, sliceName, tags, createTag, dryRun, builderConfig, happyConfig, awsBackend)
	if err != nil {
		return nil, err
	}
	workspaceRepo := createWorkspaceRepo(dryRun, awsBackend)
	stackService := stackservice.NewStackService().
		WithHappyConfig(happyConfig).
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return &HappyClient{
		HappyConfig:     happyConfig,
		StackService:    stackService,
		ArtifactBuilder: ab,
		StackTags:       stackTags,
		AWSBackend:      awsBackend,
	}, nil
}

func createWorkspaceRepo(isDryRun bool, backend *backend.Backend) workspace_repo.WorkspaceRepoIface {
	if util.IsLocalstackMode() {
		return workspace_repo.NewLocalWorkspaceRepo().WithDryRun(isDryRun)
	}
	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()
	return workspace_repo.NewWorkspaceRepo(url, org).WithDryRun(isDryRun)
}

func configureArtifactBuilder(
	ctx context.Context,
	sliceName string,
	tags []string,
	createTag, dryRun bool,
	builderConfig *ab.BuilderConfig,
	happyConfig *config.HappyConfig,
	backend *backend.Backend) (ab.ArtifactBuilderIface, map[string]string, error) {
	artifactBuilder := ab.NewArtifactBuilder(dryRun).
		WithHappyConfig(happyConfig).
		WithConfig(builderConfig).
		WithBackend(backend)
	var err error
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to find the slice %s", sliceName)
		}
		builderConfig.Profile = slice.Profile
	}

	// if creating tag and none specified, generate the default tag
	generatedTag := ""
	artifactBuilder.WithTags(tags)
	if createTag && len(artifactBuilder.GetTags()) == 0 {
		generatedTag, err = backend.GenerateTag(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to generate tag")
		}
		artifactBuilder.WithTags([]string{generatedTag})
	}

	stackTags := map[string]string{}
	if sliceName != "" {
		serviceImages, err := builderConfig.GetBuildServicesImage(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to get build service images")
		}

		for service := range serviceImages {
			stackTags[service] = generatedTag
		}
	}

	return artifactBuilder, stackTags, nil
}

type validation func() error

func validateImageExists(ctx context.Context, createTag, skipCheckTag bool, ab ab.ArtifactBuilderIface) validation {
	return func() error {
		if skipCheckTag {
			return nil
		}

		if createTag {
			// if we build and push and it succeeds, we know that the image exists
			return ab.BuildAndPush(ctx)
		}

		if len(ab.GetTags()) == 0 {
			return errors.Errorf("no tags have been assigned")
		}

		for _, tag := range ab.GetTags() {
			exists, err := ab.CheckImageExists(ctx, tag)
			if err != nil {
				return errors.Wrapf(err, "error checking if tag %s existed", tag)
			}
			if !exists {
				return errors.Errorf("image tag does not exist: '%s'", tag)
			}
		}

		return nil
	}
}
func validateTFEBackLog(ctx context.Context, isDryRun bool, awsBackend *backend.Backend) validation {
	return func() error {
		return verifyTFEBacklog(ctx, createWorkspaceRepo(isDryRun, awsBackend))
	}
}

func validateGitTree(projectRoot string) validation {
	return func() error {
		return util.ValidateGitTree(projectRoot)
	}
}

func validateStackNameAvailable(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	return func() error {
		if force {
			return nil
		}

		_, err := stackService.GetStack(ctx, stackName)
		if err != nil {
			return nil
		}

		return errors.Wrap(err, "the stack name is already taken")
	}
}

func validateConfigurationIntegirty(ctx context.Context, happyClient *HappyClient) validation {
	return func() error {
		availableServices, err := happyClient.ArtifactBuilder.GetServices(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get available services")
		}
		for _, configuredService := range happyClient.HappyConfig.GetServices() {
			if _, ok := availableServices[configuredService]; !ok {
				return errors.Errorf("service %s is not available, check your docker-compose.yml", configuredService)
			}
		}

		dir := "/Users/alokshin/GitHub/chanzuckerberg/k8s-test-app/.happy/terraform/envs/rdev"

		//TODO: Parse env tf and figure out which services are referenced

		var modules []map[string]interface{}
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".tf" {
				// Read the file contents
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// Parse the HCL syntax tree
				f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
				if diags.HasErrors() {
					return fmt.Errorf("failed to parse %s: %s", path, diags.Error())
				}
				schema := &hcl.BodySchema{
					Blocks: []hcl.BlockHeaderSchema{
						{
							Type:       "terraform",
							LabelNames: nil,
						},
						{
							Type:       "variable",
							LabelNames: []string{"name"},
						},
						{
							Type:       "output",
							LabelNames: []string{"name"},
						},
						{
							Type:       "provider",
							LabelNames: []string{"name"},
						},
						{
							Type:       "resource",
							LabelNames: []string{"type", "name"},
						},
						{
							Type:       "data",
							LabelNames: []string{"type", "name"},
						},
						{
							Type:       "module",
							LabelNames: []string{"name"},
						},
						{
							Type: "locals",
						},
					},
				}
				// Extract module blocks and their variables
				content, _, contentDiags := f.Body.PartialContent(schema)
				if contentDiags.HasErrors() {
					return errors.New("Terraform code has errors")
				}
				for _, block := range content.Blocks {
					// if block.Type == "module" {
					// 	moduleContent, _ := block.Body.Content(schema)
					// 	module := map[string]interface{}{
					// 		"name": block.Labels[0],
					// 		// "source":    moduleContent.Attributes["source"].Expr. ).Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(),
					// 		// "version":   moduleContent.GetAttribute("version").Expr.(*hclsyntax.LiteralValueExpr).Value.AsString(),
					// 		"variables": make(map[string]interface{}),
					// 	}
					// 	for _, variable := range moduleContent.Blocks {
					// 		if variable.Type == "variable" {
					// 			attributes, _ := variable.Body.JustAttributes()
					// 			module["variables"].(map[string]interface{})[variable.Labels[0]] = attributes[""] variable.Body..GetAttribute("default").Expr.String()
					// 		}
					// 	}
					// 	modules = append(modules, module)
					// }
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Print the list of modules with their variables
		for _, module := range modules {
			fmt.Printf("Module: %s\n", module["name"])
			fmt.Printf("  Source: %s\n", module["source"])
			fmt.Printf("  Version: %s\n", module["version"])
			fmt.Println("  Variables:")
			for name, value := range module["variables"].(map[string]interface{}) {
				fmt.Printf("    %s = %s\n", name, value)
			}
		}

		return nil
	}
}

func validate(validations ...validation) error {
	for _, validation := range validations {
		err := validation()
		if err != nil {
			return errors.Wrap(err, "unable to validate the environment")
		}
	}
	return nil
}

func makeWaitOptions(stackName string, happyConfig *config.HappyConfig, backend *backend.Backend) waitoptions.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(backend)
	return waitoptions.WaitOptions{
		StackName:    stackName,
		Orchestrator: taskOrchestrator,
		Services:     happyConfig.GetServices(),
	}
}

func verifyTFEBacklog(ctx context.Context, workspaceRepo workspace_repo.WorkspaceRepoIface) error {
	if !diagnostics.IsInteractiveContext(ctx) {
		// When you're not interactive, no point in measuring the backlog size
		return nil
	}
	backlogSize, _, err := workspaceRepo.EstimateBacklogSize(ctx)
	if err != nil {
		return errors.Wrap(err, "error estimating TFE backlog")
	}
	if backlogSize < 2 {
		logrus.Debug("There is no TFE backlog, proceeding.")
	} else if backlogSize < 20 {
		logrus.Debugf("TFE backlog is only %d runs long, proceeding.", backlogSize)
	} else {
		proceed := false
		prompt := &survey.Confirm{Message: fmt.Sprintf("TFE backlog is %d runs long, it might take a while to clear out. Do you want to wait? ", backlogSize)}
		err = survey.AskOne(prompt, &proceed)
		if err != nil {
			return errors.Wrapf(err, "failed to ask for confirmation")
		}

		if !proceed {
			return err
		}
	}
	return nil
}
