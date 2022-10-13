package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy-shared/model"
	cmd_util "github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stack     string
	fromEnv   string
	fromStack string
	logger    *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	rootCmd.AddCommand(configCmd)
	config.ConfigureCmdWithBootstrapConfig(configCmd)
	configCmd.PersistentFlags().StringVarP(&stack, "stack", "s", "", "Specify the stack that this applies to")

	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configDeleteCmd)

	configCmd.AddCommand(configCopyCmd)
	configCopyCmd.Flags().StringVar(&fromEnv, "from-env", "", "Specify the env that the config should be copied from")
	configCopyCmd.Flags().StringVar(&fromStack, "from-stack", "", "Specify the stack that the config should be copied from")
	configCopyCmd.MarkFlagRequired("from-env")

	configCmd.AddCommand(configDiffCmd)
	configDiffCmd.Flags().StringVar(&fromEnv, "from-env", "", "Specify the env that the config should be copied from")
	configDiffCmd.Flags().StringVar(&fromStack, "from-stack", "", "Specify the stack that the config should be copied from")
	configDiffCmd.MarkFlagRequired("from-env")
}

type ConfigRecord struct {
	Key    string `header:"Key"`
	Value  string `header:"Value"`
	Source string `header:"Source"`
}

func NewConfigRecord(record *model.ResolvedAppConfig) ConfigRecord {
	return ConfigRecord{Key: record.Key, Value: record.Value, Source: record.Source}
}

func GetHappyConfigForCmd(cmd *cobra.Command) (*config.HappyConfig, error) {
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	return happyConfig, nil
}
func ValidateConfigFeature(cmd *cobra.Command, args []string) error {
	happyConfig, err := GetHappyConfigForCmd(cmd)
	if err != nil {
		return err
	}

	if !happyConfig.GetFeatures().EnableHappyApiUsage {
		return errors.Errorf("Cannot use the %s feature set until you enable happy-api usage in your happy config json", cmd.Use)
	}

	return cmd_util.ValidateWithHappyApi(cmd, happyConfig)
}

var configCmd = &cobra.Command{
	Use:               "config",
	Short:             "modify app configs",
	Long:              "Create, Read, Update, and Delete app configs for environment '{env}'",
	SilenceUsage:      true,
	PersistentPreRunE: ValidateConfigFeature,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Usage())
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list configs",
	Long:         "List configs for the given app, env, and stack",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("listing app configs in environment '%s'", happyConfig.GetEnv()),
		))

		body := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		api := util.MakeApiClient(happyConfig)

		result := model.WrappedResolvedAppConfigsWithCount{}
		err = api.GetParsed("/v1/configs", body, &result, "attempt to list configs received 404 response")
		if err != nil {
			return err
		}

		printTable(result.Records, NewConfigRecord)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:          "get KEY",
	Short:        "get config",
	Long:         "Get the config for the given app, env, stack, and key",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("retrieving app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		api := util.MakeApiClient(happyConfig)

		result := model.WrappedResolvedAppConfig{}
		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = api.GetParsed(fmt.Sprintf("/v1/configs/%s", key), body, &result, notFoundMessage)
		if err != nil {
			return err
		}

		printTable([]*model.ResolvedAppConfig{result.Record}, NewConfigRecord)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:          "set KEY VALUE",
	Short:        "set config",
	Long:         "Set the config for the given app, env, stack, and key to the provided value",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		value := args[1]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("setting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppConfigPayload(happyConfig.App(), happyConfig.GetEnv(), stack, key, value)
		api := util.MakeApiClient(happyConfig)

		result := model.WrappedResolvedAppConfig{}
		err = api.PostParsed("/v1/configs", body, &result, "attempt to set config received 404 response")
		if err != nil {
			return err
		}

		printTable([]*model.ResolvedAppConfig{result.Record}, NewConfigRecord)
		return nil
	},
}

var configDeleteCmd = &cobra.Command{
	Use:          "delete KEY",
	Short:        "delete config",
	Long:         "Delete the config for the given app, env, stack, and key",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("deleting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppConfigLookupPayload(happyConfig.App(), happyConfig.GetEnv(), stack, key)
		api := util.MakeApiClient(happyConfig)

		result := model.WrappedResolvedAppConfig{}
		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = api.DeleteParsed(fmt.Sprintf("/v1/configs/%s", key), body, &result, notFoundMessage)
		if err != nil {
			return err
		}

		if result.Record != nil {
			logrus.Infof("app config with key '%s' has been deleted", result.Record.Key)
		} else {
			return errors.New(notFoundMessage)
		}

		return nil
	},
}

var configCopyCmd = &cobra.Command{
	Use:          "cp KEY",
	Short:        "copy config",
	Long:         "Copy the config for the given app, source env, source stack, and key to the given destination env and destination stack",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		srcAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("copying app config with key '%s' from %s to %s", key, srcAppEnvStack, destAppEnvStack)

		body := model.NewAppConfigDiffPayload(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack)
		api := util.MakeApiClient(happyConfig)

		result := model.WrappedResolvedAppConfig{}
		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = api.GetParsed("/v1/config/diff", body, &result, notFoundMessage)
		if err != nil {
			return err
		}

		if result.Record != nil {
			logrus.Infof("app config with key '%s' has been copied from %s to %s", result.Record.Key, srcAppEnvStack, destAppEnvStack)
		} else {
			return errors.New(notFoundMessage)
		}

		return nil
	},
}

var configDiffCmd = &cobra.Command{
	Use:          "diff",
	Short:        "diff config",
	Long:         "Get a list of config keys that are present in the given app, source env, source stack but not in the given destination env and destination stack",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		srcAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("retrieving list of config keys that exist in %s and not %s", srcAppEnvStack, destAppEnvStack)

		body := model.NewAppConfigDiffPayload(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack)
		api := util.MakeApiClient(happyConfig)

		result := model.ConfigDiffResponse{}
		err = api.GetParsed("/v1/config/diff", body, &result, "attempt to get config diff received 404 response")
		if err != nil {
			return err
		}

		logrus.Infof("the following keys are present in %s and not in %s", srcAppEnvStack, destAppEnvStack)
		tablePrinter := util.NewTablePrinter()
		tablePrinter.Print(result.MissingKeys)
		return nil
	},
}

func printTable[T interface{}, Z interface{}](rows []T, rowStruct func(record T) Z) {
	tablePrinter := util.NewTablePrinter()
	for _, row := range rows {
		tablePrinter.AddRow(rowStruct(row))
	}
	tablePrinter.Flush()
}

func messageWithStackSuffix(message string) string {
	stackSuffix := ""
	if stack != "" {
		stackSuffix = ", stack '" + stack + "'"
	}
	return message + stackSuffix
}
