package cmd

import (
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sourceTag string
var destTags []string
var images []string

func init() {
	rootCmd.AddCommand(tagsCmd)
	config.ConfigureCmdWithBootstrapConfig(tagsCmd)

	tagsCmd.Flags().StringVar(&sourceTag, "source-tag", "", "Tag name for existing docker image.")
	tagsCmd.Flags().StringSliceVar(&destTags, "dest-tag", []string{}, "Extra tags to apply and push to the docker repo.")
	_ = tagsCmd.MarkFlagRequired("source-tag")
	_ = tagsCmd.MarkFlagRequired("dest-tag")
}

var tagsCmd = &cobra.Command{
	Use:          "addtags <stack-name>",
	Short:        "Add additional tags to already-pushed images in the ECR repo",
	Long:         "Add additional tags to already-pushed images in the ECR repo",
	SilenceUsage: true,
	RunE:         runTags,
	PreRunE: happyCmd.Validate(
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(),
				checklist.AwsInstalled,
			)
		},
	),
}

func runTags(cmd *cobra.Command, args []string) error {
	stackName := args[0]
	happyClient, err := makeHappyClient(cmd, sliceName, stackName, tags, createTag)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}
	serviceRegistries := happyClient.AWSBackend.GetIntegrationSecret().GetServiceRegistries()
	stackECRS, err := happyClient.ArtifactBuilder.GetECRsForServices(cmd.Context())
	if err != nil {
		log.Debugf("unable to get ECRs for services: %s", err)
	}
	if len(stackECRS) > 0 {
		serviceRegistries = stackECRS
	}
	return happyClient.ArtifactBuilder.RetagImages(cmd.Context(), serviceRegistries, sourceTag, destTags, images)
}
