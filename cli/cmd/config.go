package cmd

import (
	"fmt"

	cmd_util "github.com/chanzuckerberg/happy/pkg/cli/cmd"
	"github.com/chanzuckerberg/happy/pkg/cli/config"
	"github.com/chanzuckerberg/happy/pkg/cli/util"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/shared/model"
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
}

type ConfigTableEntry struct {
	Key    string `header:"Key"`
	Value  string `header:"Value"`
	Source string `header:"Source"`
}

func newConfigTableEntry(record *model.ResolvedAppConfig) ConfigTableEntry {
	return ConfigTableEntry{Key: record.Key, Value: record.Value, Source: record.Source}
}

func ValidateConfigFeature(cmd *cobra.Command, args []string) error {
	happyConfig, err := config.GetHappyConfigForCmd(cmd)
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
		logrus.Println(cmd.Usage())
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list configs",
	Long:         "List configs for the given app, env, and stack",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("listing app configs in environment '%s'", happyConfig.GetEnv()),
		))

		api := util.MakeApiClient(happyConfig)
		result, err := api.ListConfigs(happyConfig.App(), happyConfig.GetEnv(), stack)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New("attempt to list configs received 404 response")
			}
			return err
		}

		printTable(result.Records, newConfigTableEntry)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:          "get KEY",
	Short:        "get config",
	Long:         "Get the config for the given app, env, stack, and key",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("retrieving app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)

		api := util.MakeApiClient(happyConfig)
		result, err := api.GetConfig(happyConfig.App(), happyConfig.GetEnv(), stack, key)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New(notFoundMessage)
			}
			return err
		}

		printTable([]*model.ResolvedAppConfig{result.Record}, newConfigTableEntry)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:          "set KEY VALUE",
	Short:        "set config",
	Long:         "Set the config for the given app, env, stack, and key to the provided value",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		value := args[1]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("setting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		api := util.MakeApiClient(happyConfig)
		result, err := api.SetConfig(happyConfig.App(), happyConfig.GetEnv(), stack, key, value)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New("attempt to set config received 404 response")
			}
			return err
		}

		printTable([]*model.ResolvedAppConfig{result.Record}, newConfigTableEntry)
		return nil
	},
}

var configDeleteCmd = &cobra.Command{
	Use:          "delete KEY",
	Short:        "delete config",
	Long:         "Delete the config for the given app, env, stack, and key",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("deleting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)

		api := util.MakeApiClient(happyConfig)
		result, err := api.DeleteConfig(happyConfig.App(), happyConfig.GetEnv(), stack, key)
		if err != nil && !errors.Is(err, client.ErrRecordNotFound) {
			return err
		}

		if result.Record == nil {
			return errors.New(notFoundMessage)
		}

		logrus.Infof("app config with key '%s' has been deleted", result.Record.Key)
		return nil
	},
}

var configCopyCmd = &cobra.Command{
	Use:          "cp KEY",
	Short:        "copy config",
	Long:         "Copy the config for the given app, source env, source stack, and key to the given destination env and destination stack",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		srcAppEnvStack := model.NewAppMetadata(happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("copying app config with key '%s' from %s to %s", key, srcAppEnvStack, destAppEnvStack)

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)

		api := util.MakeApiClient(happyConfig)
		result, err := api.CopyConfig(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack, key)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		srcAppEnvStack := model.NewAppMetadata(happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("retrieving list of config keys that exist in %s and not %s", srcAppEnvStack, destAppEnvStack)

		api := util.MakeApiClient(happyConfig)
		result, err := api.GetMissingConfigKeys(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack)
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
		stackSuffix = fmt.Sprintf(", stack '%s'", stack)
	}
	return message + stackSuffix
}
