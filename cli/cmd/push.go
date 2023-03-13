package cmd

import (
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/shared/opts"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var tags []string
var sliceName string

func init() {
	rootCmd.AddCommand(pushCmd)
	config.ConfigureCmdWithBootstrapConfig(pushCmd)

	pushCmd.Flags().StringVar(&sliceName, "slice", "", "The name of the slice you'd like to push to the registry.")
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
		happyCmd.IsStackNameAlphaNumeric),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		happyClient, err := makeHappyClient(cmd, sliceName, stackName, tag, createTag, dryRun)
		if err != nil {
			return errors.Wrap(err, "unable to initialize the happy client")
		}

		ctx := cmd.Context()
		dryRunOption := opts.DryRun(dryRun)
		err = validate(
			validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
			validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
			validateStackExistsCreate(ctx, stackName, happyClient, dryRunOption),
			validateECRExists(ctx, stackName, terraformECRTargetPathTemplate, happyClient, dryRunOption),
		)
		if err != nil {
			return errors.Wrap(err, "failed one of the happy client validations")
		}

		return happyClient.ArtifactBuilder.BuildAndPush(ctx)
	},
}
