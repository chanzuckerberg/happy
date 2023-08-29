package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	force         bool
	skipCheckTag  bool
	createTag     bool
	tag           string
	dryRun        bool
	imageSrcEnv   string
	imageSrcStack string
)

func init() {
	rootCmd.AddCommand(createCmd)
	config.ConfigureCmdWithBootstrapConfig(createCmd)
	happyCmd.SupportUpdateSlices(createCmd, &sliceName, &sliceDefaultTag) // Should this function be renamed to something more generalized?
	happyCmd.SetMigrationFlags(createCmd)
	happyCmd.SetImagePromotionFlags(createCmd, &imageSrcEnv, &imageSrcStack)
	happyCmd.SetDryRunFlag(createCmd, &dryRun)
	createCmd.Flags().StringVar(&tag, "tag", "", "Specify the tag for the docker images. If not specified we will generate a default tag.")
	createCmd.Flags().BoolVar(&createTag, "create-tag", true, "Will build, tag, and push images when set. Otherwise, assumes images already exist.")
	createCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
	createCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
}

var createCmd = &cobra.Command{
	Use:          "create STACK_NAME",
	Short:        "Create new stack",
	Long:         "Create a new stack with a given tag.",
	SilenceUsage: true,
	PreRunE: happyCmd.Validate(
		happyCmd.IsImageEnvUsedWithImageStack,
		happyCmd.IsTagUsedWithSkipTag,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()

			required_checks := []util.ValidationCallback{
				checklist.TerraformInstalled,
				checklist.AwsInstalled,
			}

			if !skipCheckTag || createTag {
				required_checks = append(required_checks, checklist.MinDockerComposeVersion, checklist.DockerEngineRunning, checklist.DockerInstalled)
			}

			return util.ValidateEnvironment(cmd.Context(), required_checks...)
		},
	),
	RunE: runCreate,
}

// keep in sync with happy-stack-eks terraform module
const terraformECRTargetPathTemplate = `module.stack.module.services["%s"].module.ecr`

func runCreate(
	cmd *cobra.Command,
	args []string,
) (err error) {
	stackName := args[0]

	happyClient, err := makeHappyClient(cmd, sliceName, stackName, []string{tag}, createTag)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}
	ctx := context.WithValue(cmd.Context(), options.DryRunKey, dryRun)
	message := workspace_repo.Message(fmt.Sprintf("Happy %s Create Stack [%s]", util.GetVersion().Version, stackName))
	err = validate(
		validateConfigurationIntegirty(ctx, sliceName, happyClient),
		validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
		validateStackNameGloballyAvailable(ctx, happyClient.StackService, stackName, force),
		validateTFEBackLog(ctx, happyClient.AWSBackend),
		validateStackExistsCreate(ctx, stackName, happyClient, message),
		validateECRExists(ctx, stackName, terraformECRTargetPathTemplate, happyClient, message),
		validateImageExists(ctx, createTag, skipCheckTag, imageSrcEnv, imageSrcStack, happyClient),
	)
	if err != nil {
		return errors.Wrap(err, "failed one of the happy client validations")
	}

	// update the newly created stack
	stack, err := happyClient.StackService.GetStack(ctx, stackName)
	if err != nil {
		return errors.Wrapf(err, "stack %s doesn't exist; this should never happen", stackName)
	}

	err = updateStack(ctx, cmd, stack, force, happyClient)
	if err != nil {
		return errors.Wrapf(err, "unable to update the stack %s", stack.Name)
	}
	// if it was a dry run, we should remove the stack after we are done
	if dryRun {
		log.Debugf("cleaning up stack '%s'", stack.Name)
		return errors.Wrap(happyClient.StackService.Remove(ctx, stack.Name), "unable to remove stack")
	}
	return nil
}

func validateECRExists(ctx context.Context, stackName string, ecrTargetPathFormat string, happyClient *HappyClient, options ...workspace_repo.TFERunOption) validation {
	log.Debug("Scheduling validateECRExists()")
	return func() error {
		log.Debug("Running validateECRExists()")
		if !happyClient.HappyConfig.GetFeatures().EnableECRAutoCreation {
			return nil
		}

		stackECRS, err := happyClient.ArtifactBuilder.GetECRsForServices(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get ECRs for services; this shouldn't happen if the stack TF is configured correctly")
		}

		missingServiceECRs := []string{}
		for _, service := range happyClient.HappyConfig.GetServices() {
			if _, ok := stackECRS[service]; !ok {
				missingServiceECRs = append(missingServiceECRs, service)
			}
		}
		if len(missingServiceECRs) == 0 {
			return nil
		}

		log.Debugf("missing ECRs for the following services %s. making them now", strings.Join(missingServiceECRs, ","))
		targetAddrs := []string{}
		for _, service := range happyClient.HappyConfig.GetServices() {
			targetAddrs = append(targetAddrs, fmt.Sprintf(ecrTargetPathFormat, service))
		}
		stack, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "stack %s doesn't exist; this should never happen", stackName)
		}
		stackMeta, err := updateStackMeta(ctx, stack.Name, happyClient)
		if err != nil {
			return errors.Wrap(err, "unable to update the stack's meta information")
		}

		// this has a strong coupling with the TF module version that we are using in happy-stack-eks,
		// so if the user isn't on it yet, this will fail or not do what you are expecting
		// TODO: maybe CDK
		// TODO: maybe we can peek at the version and fail if its not right or something?
		stack = stack.WithMeta(stackMeta)
		tfDirPath := happyClient.HappyConfig.TerraformDirectory()
		happyProjectRoot := happyClient.HappyConfig.GetProjectRoot()
		srcDir := filepath.Join(happyProjectRoot, tfDirPath)
		return stack.Apply(ctx, srcDir, makeWaitOptions(stackName, happyClient.HappyConfig, happyClient.AWSBackend), append(options, workspace_repo.TargetAddrs(targetAddrs))...)
	}
}

func validateStackExistsCreate(ctx context.Context, stackName string, happyClient *HappyClient, options ...workspace_repo.TFERunOption) validation {
	log.Debug("Scheduling validateStackExistsCreate()")
	return func() error {
		log.Debug("Running validateStackExistsCreate()")
		// 1.) if the stack does not exist and force flag is used, call the create function first
		_, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			log.Debugf("Stack doesn't exist %s: %s\n", stackName, err.Error())
			_, err = happyClient.StackService.Add(ctx, stackName, options...)
			if err != nil {
				return errors.Wrap(err, "unable to create the stack")
			}
			log.Debugf("Stack added: %s", stackName)
		} else {
			log.Debugf("Stack exists: %s", stackName)
			if !force {
				return errors.Wrapf(err, "stack %s already exists", stackName)
			}
		}

		return nil
	}
}
