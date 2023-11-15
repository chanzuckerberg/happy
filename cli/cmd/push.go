package cmd

import (
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var tags []string
var sliceName string

func init() {
	rootCmd.AddCommand(pushCmd)
	config.ConfigureCmdWithBootstrapConfig(pushCmd)

	pushCmd.Flags().StringVar(&sliceName, "slice", "", "The name of the slice you'd like to push to the registry.")
	pushCmd.Flags().StringSliceVar(&tags, "tags", nil, "Extra tags to set for built images, comma-delimited (ex: tag1,tag2,tag3). We will, in addition, generate default tags automatically.")
	pushCmd.Flags().StringSliceVar(&tags, "tag", nil, "Extra tags to set for built images, comma-delimited (ex: tag1,tag2,tag3). We will, in addition, generate default tags automatically.")
}

var pushCmd = &cobra.Command{
	Use:          "push STACK_NAME",
	Short:        "Push docker images",
	Long:         "Push docker images to ECR",
	SilenceUsage: true,
	PreRunE: happyCmd.Validate(
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(),
				checklist.DockerEngineRunning,
				checklist.MinDockerComposeVersion,
				checklist.DockerInstalled,
				checklist.TerraformInstalled,
				checklist.AwsInstalled,
			)
		},
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		happyClient, err := makeHappyClient(cmd, sliceName, stackName, tags, createTag)
		if err != nil {
			return errors.Wrap(err, "unable to initialize the happy client")
		}

		ctx := cmd.Context()
		err = validate(
			validateConfigurationIntegirty(ctx, sliceName, happyClient),
			validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
			validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
			validateStackExistsCreate(ctx, stackName, happyClient),
			validateECRExists(ctx, stackName, happyClient),
		)
		if err != nil {
			return errors.Wrap(err, "failed one of the happy client validations")
		}

		return happyClient.ArtifactBuilder.BuildAndPush(ctx)
	},
}
