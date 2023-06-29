package cmd

import (
	"context"
	"io"

	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/cli/pkg/output"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	listAll bool
	remote  bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
	listCmd.Flags().StringVar(&OutputFormat, "output", "text", "Output format. One of: json, yaml, or text. Defaults to text, which is the only interactive mode.")
	listCmd.Flags().BoolVar(&listAll, "all", false, "List all stacks, not just those belonging to this app")
	listCmd.Flags().BoolVar(&remote, "remote", false, "List stacks from the remote happy server")
}

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List stacks",
	Long:         "Listing stacks in environment '{env}'",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if OutputFormat != "text" {
			logrus.SetOutput(io.Discard)
		}
		happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
		if err != nil {
			return errors.Wrap(err, "unable to initialize the happy client")
		}

		var metas []*model.AppStackResponse
		if remote {
			metas, err = listStacksRemote(cmd.Context(), listAll, happyClient)
			if err != nil {
				return err
			}
		} else {
			m, err := happyClient.StackService.CollectStackInfo(cmd.Context(), happyClient.HappyConfig.App())
			if err != nil {
				return errors.Wrap(err, "unable to collect stack info")
			}

			for _, meta := range m {
				// only show the stacks that belong to this app or they want to list all
				if listAll || (meta.AppMetadata.App.AppName == happyClient.HappyConfig.App()) {
					metas = append(metas, meta)
				}
			}
		}
		printer := output.NewPrinter(OutputFormat)
		err = printer.PrintStacks(cmd.Context(), metas)
		if err != nil {
			return errors.Wrap(err, "unable to print stacks")
		}

		return nil
	},
}

func listStacksRemote(ctx context.Context, listAll bool, happyClient *HappyClient) ([]*model.AppStackResponse, error) {
	api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
	result, err := api.ListStacks(model.MakeAppStackPayload(
		happyClient.HappyConfig.App(),
		happyClient.HappyConfig.GetEnv(),
		"", model.AWSContext{
			AWSProfile:     *happyClient.HappyConfig.AwsProfile(),
			AWSRegion:      *happyClient.HappyConfig.AwsRegion(),
			TaskLaunchType: "k8s",
			K8SNamespace:   happyClient.HappyConfig.K8SConfig().Namespace,
			K8SClusterID:   happyClient.HappyConfig.K8SConfig().ClusterID,
		},
	))
	if err != nil {
		return nil, err
	}

	metas := []*model.AppStackResponse{}
	for _, meta := range result.Records {
		// only show the stacks that belong to this app or they want to list all
		if listAll || (meta.AppMetadata.App.AppName == happyClient.HappyConfig.App()) {
			metas = append(metas, meta)
		}
	}

	return metas, nil
}
