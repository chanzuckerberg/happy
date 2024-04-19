package cmd

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"sort"

	cmd_util "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/shared/config"
	apiclient "github.com/chanzuckerberg/happy/shared/hapi"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stack     string
	fromEnv   string
	fromStack string
	reveal    bool
	v2        bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	config.ConfigureCmdWithBootstrapConfig(configCmd)
	configCmd.PersistentFlags().StringVarP(&stack, "stack", "s", "", "Specify the stack that this applies to")
	configCmd.PersistentFlags().BoolVar(&v2, "v2", false, "Use v2 configs (stored in K8s secrets)")

	configCmd.AddCommand(configListCmd)
	configListCmd.Flags().BoolVarP(&reveal, "reveal", "r", false, "Print the actual app config values instead of masking them")

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configDeleteCmd)

	configCmd.AddCommand(configCopyCmd)
	configCopyCmd.Flags().StringVar(&fromEnv, "from-env", "", "Specify the env that the config should be copied from")
	configCopyCmd.Flags().StringVar(&fromStack, "from-stack", "", "Specify the stack that the config should be copied from")
	err := configCopyCmd.MarkFlagRequired("from-env")
	if err != nil {
		logrus.Panic("failed to mark flag as required")
	}

	configCmd.AddCommand(configDiffCmd)
	configDiffCmd.Flags().StringVar(&fromEnv, "from-env", "", "Specify the env that the config should be copied from")
	configDiffCmd.Flags().StringVar(&fromStack, "from-stack", "", "Specify the stack that the config should be copied from")
	err = configDiffCmd.MarkFlagRequired("from-env")
	if err != nil {
		logrus.Panic("failed to mark flag as required")
	}

	configCmd.AddCommand(configExecCmd)
}

type ConfigTableEntry struct {
	Key    string `header:"Key"`
	Value  string `header:"Value"`
	Source string `header:"Source"`
}

func newConfigTableEntry(record *model.ResolvedAppConfig) ConfigTableEntry {
	return ConfigTableEntry{Key: record.Key, Value: record.Value, Source: record.Source}
}

func newConfigTableEntryV2(record apiclient.AppConfigList) ConfigTableEntry {
	return ConfigTableEntry{Key: record.Key, Value: record.Value, Source: string(record.Source)}
}

func ValidateConfigFeature(cmd *cobra.Command, args []string) error {
	happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
	if err != nil {
		return err
	}

	if !happyClient.HappyConfig.GetFeatures().EnableHappyApiUsage {
		return errors.Errorf("Cannot use the %s feature set until you enable happy-api usage in your happy config json", cmd.Use)
	}

	return cmd_util.ValidateWithHappyApi(cmd, happyClient.HappyConfig, happyClient.AWSBackend)
}

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "modify app configs",
	Long:         "Create, Read, Update, and Delete app configs for environment '{env}'",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := ValidateConfigFeature(cmd, args)
		if err != nil {
			return err
		}

		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.DockerEngineRunning,
			checklist.MinDockerComposeVersion,
			checklist.DockerInstalled,
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Println(cmd.Usage())
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list configs",
	Long:         "List configs for the given app, env, and stack",
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

		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("listing app configs in environment '%s'", happyClient.HappyConfig.GetEnv()),
		))

		if happyClient.HappyConfig.GetFeatures().EnableHappyConfigV2 || v2 {
			awsCredsProvider := hapi.NewAWSCredentialsProviderCLI(happyClient.AWSBackend)
			creds, err := awsCredsProvider.GetCredentials(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to get aws credentials")
			}

			params := &apiclient.ListAppConfigParams{
				AppName:             happyClient.HappyConfig.App(),
				Environment:         happyClient.HappyConfig.GetEnv(),
				Stack:               &stack,
				AwsProfile:          *happyClient.HappyConfig.AwsProfile(),
				AwsRegion:           *happyClient.HappyConfig.AwsRegion(),
				K8sNamespace:        happyClient.HappyConfig.K8SConfig().Namespace,
				K8sClusterId:        happyClient.HappyConfig.K8SConfig().ClusterID,
				XAwsAccessKeyId:     b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				XAwsSecretAccessKey: b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
				XAwsSessionToken:    creds.SessionToken, // SessionToken is already base64 encoded
			}
			api := hapi.MakeAPIClientV2(happyClient.HappyConfig)
			resp, err := api.ListAppConfigWithResponse(context.Background(), params)
			if err != nil {
				return err
			}
			if resp.HTTPResponse.StatusCode >= 400 {
				return errors.New(string(resp.Body))
			}

			printTableV2(*resp.JSON200, newConfigTableEntryV2, !reveal)
		} else {
			api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
			result, err := api.ListConfigs(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack)
			if err != nil {
				if errors.Is(err, client.ErrRecordNotFound) {
					return errors.New("attempt to list configs received 404 response")
				}
				return err
			}

			printTable(result.Records, newConfigTableEntry, !reveal)
		}
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:          "get KEY",
	Short:        "get config",
	Long:         "Get the config for the given app, env, stack, and key",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
			fmt.Sprintf("retrieving app config with key '%s' in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		))

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		)

		if happyClient.HappyConfig.GetFeatures().EnableHappyConfigV2 || v2 {
			awsCredsProvider := hapi.NewAWSCredentialsProviderCLI(happyClient.AWSBackend)
			creds, err := awsCredsProvider.GetCredentials(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to get aws credentials")
			}

			params := &apiclient.ReadAppConfigParams{
				AppName:             happyClient.HappyConfig.App(),
				Environment:         happyClient.HappyConfig.GetEnv(),
				Stack:               &stack,
				AwsProfile:          *happyClient.HappyConfig.AwsProfile(),
				AwsRegion:           *happyClient.HappyConfig.AwsRegion(),
				K8sNamespace:        happyClient.HappyConfig.K8SConfig().Namespace,
				K8sClusterId:        happyClient.HappyConfig.K8SConfig().ClusterID,
				XAwsAccessKeyId:     b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				XAwsSecretAccessKey: b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
				XAwsSessionToken:    creds.SessionToken, // SessionToken is already base64 encoded
			}
			api := hapi.MakeAPIClientV2(happyClient.HappyConfig)
			resp, err := api.ReadAppConfigWithResponse(context.Background(), key, params)
			if err != nil {
				return err
			}
			if resp.HTTPResponse.StatusCode >= 400 {
				return errors.New(string(resp.Body))
			}

			printTableV2([]apiclient.AppConfigList{*resp.JSON200}, newConfigTableEntryV2, false)
		} else {
			api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
			result, err := api.GetConfig(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack, key)
			if err != nil {
				if errors.Is(err, client.ErrRecordNotFound) {
					return errors.New(notFoundMessage)
				}
				return err
			}
			printTable([]*model.ResolvedAppConfig{result.Record}, newConfigTableEntry, false)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:          "set KEY VALUE",
	Short:        "set config",
	Long:         "Set the config for the given app, env, stack, and key to the provided value",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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

		key := args[0]
		value := args[1]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("setting app config with key '%s' in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		))

		if happyClient.HappyConfig.GetFeatures().EnableHappyConfigV2 || v2 {
			awsCredsProvider := hapi.NewAWSCredentialsProviderCLI(happyClient.AWSBackend)
			creds, err := awsCredsProvider.GetCredentials(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to get aws credentials")
			}

			params := &apiclient.SetAppConfigParams{
				AppName:             happyClient.HappyConfig.App(),
				Environment:         happyClient.HappyConfig.GetEnv(),
				Stack:               &stack,
				AwsProfile:          *happyClient.HappyConfig.AwsProfile(),
				AwsRegion:           *happyClient.HappyConfig.AwsRegion(),
				K8sNamespace:        happyClient.HappyConfig.K8SConfig().Namespace,
				K8sClusterId:        happyClient.HappyConfig.K8SConfig().ClusterID,
				XAwsAccessKeyId:     b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				XAwsSecretAccessKey: b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
				XAwsSessionToken:    creds.SessionToken, // SessionToken is already base64 encoded
			}
			api := hapi.MakeAPIClientV2(happyClient.HappyConfig)
			resp, err := api.SetAppConfigWithResponse(context.Background(), params, apiclient.SetAppConfigJSONRequestBody{Key: key, Value: value})
			if err != nil {
				return err
			}
			if resp.HTTPResponse.StatusCode >= 400 {
				return errors.New(string(resp.Body))
			}
			key = resp.JSON200.Key
		} else {
			api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
			resp, err := api.SetConfig(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack, key, value)
			if err != nil {
				if errors.Is(err, client.ErrRecordNotFound) {
					return errors.New("attempt to set config received 404 response")
				}
				return err
			}
			key = resp.Record.Key
		}
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("successfully set app config with key '%s' in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		))
		return nil
	},
}

var configDeleteCmd = &cobra.Command{
	Use:          "delete KEY",
	Short:        "delete config",
	Long:         "Delete the config for the given app, env, stack, and key",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("deleting app config with key '%s' in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		))

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		)

		if happyClient.HappyConfig.GetFeatures().EnableHappyConfigV2 || v2 {
			awsCredsProvider := hapi.NewAWSCredentialsProviderCLI(happyClient.AWSBackend)
			creds, err := awsCredsProvider.GetCredentials(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to get aws credentials")
			}

			params := &apiclient.DeleteAppConfigParams{
				AppName:             happyClient.HappyConfig.App(),
				Environment:         happyClient.HappyConfig.GetEnv(),
				Stack:               &stack,
				AwsProfile:          *happyClient.HappyConfig.AwsProfile(),
				AwsRegion:           *happyClient.HappyConfig.AwsRegion(),
				K8sNamespace:        happyClient.HappyConfig.K8SConfig().Namespace,
				K8sClusterId:        happyClient.HappyConfig.K8SConfig().ClusterID,
				XAwsAccessKeyId:     b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				XAwsSecretAccessKey: b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
				XAwsSessionToken:    creds.SessionToken, // SessionToken is already base64 encoded
			}
			api := hapi.MakeAPIClientV2(happyClient.HappyConfig)
			resp, err := api.DeleteAppConfigWithResponse(context.Background(), key, params)
			if err != nil {
				return err
			}
			if resp.HTTPResponse.StatusCode >= 400 {
				return errors.New(string(resp.Body))
			}
		} else {
			api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
			result, err := api.DeleteConfig(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack, key)
			if err != nil && !errors.Is(err, client.ErrRecordNotFound) {
				return err
			}

			if result.Record == nil {
				return errors.New(notFoundMessage)
			}
		}

		logrus.Infof("app config with key '%s' has been deleted", key)
		return nil
	},
}

var configCopyCmd = &cobra.Command{
	Use:          "cp KEY",
	Short:        "copy config",
	Long:         "Copy the config for the given app, source env, source stack, and key to the given destination env and destination stack",
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

		key := args[0]
		srcAppEnvStack := model.NewAppMetadata(happyClient.HappyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := model.NewAppMetadata(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack)
		logrus.Infof("copying app config with key '%s' from %s to %s", key, srcAppEnvStack, destAppEnvStack)

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyClient.HappyConfig.GetEnv()),
		)

		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
		result, err := api.CopyConfig(happyClient.HappyConfig.App(), fromEnv, fromStack, happyClient.HappyConfig.GetEnv(), stack, key)
		if err != nil && !errors.Is(err, client.ErrRecordNotFound) {
			return err
		}

		if result.Record == nil {
			return errors.New(notFoundMessage)
		}

		logrus.Infof("app config with key '%s' has been copied from %s to %s", result.Record.Key, srcAppEnvStack, destAppEnvStack)
		return nil
	},
}

var configDiffCmd = &cobra.Command{
	Use:          "diff",
	Short:        "diff config",
	Long:         "Get a list of config keys that are present in the given app, source env, source stack but not in the given destination env and destination stack",
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

		srcAppEnvStack := model.NewAppMetadata(happyClient.HappyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := model.NewAppMetadata(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack)
		logrus.Infof("retrieving list of config keys that exist in %s and not %s", srcAppEnvStack, destAppEnvStack)

		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
		result, err := api.GetMissingConfigKeys(happyClient.HappyConfig.App(), fromEnv, fromStack, happyClient.HappyConfig.GetEnv(), stack)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New("attempt to get config diff received 404 response")
			}
			return err
		}

		if len(result.MissingKeys) == 0 {
			logrus.Infof("there are no config keys present in %s and not in %s", srcAppEnvStack, destAppEnvStack)
		} else {
			logrus.Infof("the following keys are present in %s and not in %s", srcAppEnvStack, destAppEnvStack)
			tablePrinter := util.NewTablePrinter()
			tablePrinter.Print(result.MissingKeys)
		}
		return nil
	},
}

func printTable[Z interface{}](rows []*model.ResolvedAppConfig, rowStruct func(record *model.ResolvedAppConfig) Z, maskValue bool) {
	tablePrinter := util.NewTablePrinter()
	for _, row := range sortAppConfigRows(rows) {
		if maskValue {
			row.Value = "********"
		}
		tablePrinter.AddRow(rowStruct(row))
	}
	tablePrinter.Flush()
}

func sortAppConfigRows(rows []*model.ResolvedAppConfig) []*model.ResolvedAppConfig {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key < rows[j].Key
	})
	return rows
}

func messageWithStackSuffix(message string) string {
	stackSuffix := ""
	if stack != "" {
		stackSuffix = fmt.Sprintf(", stack '%s'", stack)
	}
	return message + stackSuffix
}

func printTableV2[Z interface{}](rows []apiclient.AppConfigList, rowStruct func(record apiclient.AppConfigList) Z, maskValue bool) {
	tablePrinter := util.NewTablePrinter()
	for _, row := range sortAppConfigRowsV2(rows) {
		if maskValue {
			row.Value = "********"
		}
		tablePrinter.AddRow(rowStruct(row))
	}
	tablePrinter.Flush()
}

func sortAppConfigRowsV2(rows []apiclient.AppConfigList) []apiclient.AppConfigList {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key < rows[j].Key
	})
	return rows
}
