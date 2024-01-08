package cmd

import (
	"context"
	b64 "encoding/base64"
	"fmt"

	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/client"
	apiclient "github.com/chanzuckerberg/happy/shared/hapi"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configMigrateCmd)
	configCmd.AddCommand(configMigrateAllCmd)
}

var configMigrateCmd = &cobra.Command{
	Use:          "migrate KEY",
	Short:        "migrate config value from v1 to v2",
	Long:         "Migrate the config for the given app, env, stack, and key from v1 to v2",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			noKeyProvidedMessage := messageWithStackSuffix("Please supply the key name you want to look up.")
			return errors.New(noKeyProvidedMessage)
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("migrating app config with key '%s' in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		))

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("[v1] app config with key '%s' could not be found in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		)

		// Get the value from the v1 api
		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
		result, err := api.GetConfig(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack, key)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New(notFoundMessage)
			}
			return err
		}

		return migrateConfig(happyClient, result.Record)
	},
}

var configMigrateAllCmd = &cobra.Command{
	Use:          "migrate-all",
	Short:        "migrate all config values from v1 to v2",
	Long:         "Migrate all configs for the given app, env, and stack from v1 to v2",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
		if err != nil {
			return err
		}

		// Get the values from the v1 api
		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
		results, err := api.ListConfigs(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack)
		if err != nil {
			return err
		}

		for _, record := range results.Records {
			logrus.Info(messageWithStackSuffix(
				fmt.Sprintf("migrating app config with key '%s' in environment '%s'", record.Key, happyClient.HappyConfig.GetEnv()),
			))

			err := migrateConfig(happyClient, record)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func migrateConfig(happyClient *HappyClient, record *model.ResolvedAppConfig) error {
	// Write the value to the v2 api
	awsCredsProvider := hapi.NewAWSCredentialsProviderCLI(happyClient.AWSBackend)
	creds, err := awsCredsProvider.GetCredentials(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get aws credentials")
	}

	params := &apiclient.SetAppConfigParams{
		AppName:             happyClient.HappyConfig.App(),
		Environment:         happyClient.HappyConfig.GetEnv(),
		Stack:               &record.Stack,
		AwsProfile:          *happyClient.HappyConfig.AwsProfile(),
		AwsRegion:           *happyClient.HappyConfig.AwsRegion(),
		K8sNamespace:        happyClient.HappyConfig.K8SConfig().Namespace,
		K8sClusterId:        happyClient.HappyConfig.K8SConfig().ClusterID,
		XAwsAccessKeyId:     b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
		XAwsSecretAccessKey: b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
		XAwsSessionToken:    creds.SessionToken, // SessionToken is already base64 encoded
	}
	apiv2 := hapi.MakeAPIClientV2(happyClient.HappyConfig)
	resp, err := apiv2.SetAppConfigWithResponse(context.Background(), params, apiclient.SetAppConfigJSONRequestBody{Key: record.Key, Value: record.Value})
	if err != nil {
		return err
	}
	if resp.HTTPResponse.StatusCode >= 400 {
		return errors.New(string(resp.Body))
	}

	logrus.Info(messageWithStackSuffix(
		fmt.Sprintf("successfully migrated app config with key '%s' in environment '%s'", resp.JSON200.Key, happyClient.HappyConfig.GetEnv()),
	))
	return nil
}
