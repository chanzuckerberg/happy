package cmd

import (
	"fmt"
	"net/http"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
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

func NewConfigRecord(record model.ResolvedAppConfig) ConfigRecord {
	return ConfigRecord{Key: record.Key, Value: record.Value, Source: record.Source}
}

var configCmd = &cobra.Command{
	Use:               "config",
	Short:             "modify app configs",
	Long:              "Create, Read, Update, and Delete app configs for environment '{env}'",
	SilenceUsage:      true,
	PersistentPreRunE: api.ValidateConfigFeature,
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
		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("listing app configs in environment '%s'", happyConfig.GetEnv()),
		))

		body := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Get("/v1/configs", body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		err = inspectForErrors(resp, "attempt to list configs received 404 response")
		if err != nil {
			return err
		}

		result := model.WrappedResolvedAppConfigsWithCount{}
		api.ParseResponse(resp, &result)
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
		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("retrieving app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppMetadata(happyConfig.App(), happyConfig.GetEnv(), stack)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Get(fmt.Sprintf("/v1/configs/%s", key), body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = inspectForErrors(resp, notFoundMessage)
		if err != nil {
			return err
		}

		result := model.WrappedResolvedAppConfig{}
		api.ParseResponse(resp, &result)
		printTable([]model.ResolvedAppConfig{*result.Record}, NewConfigRecord)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:          "set KEY VALUE",
	Short:        "set config",
	Long:         "Set the config for the given app, env, stack, and key to the provided value",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		value := args[1]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("setting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppConfigPayload(happyConfig.App(), happyConfig.GetEnv(), stack, key, value)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Post("/v1/configs", body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		err = inspectForErrors(resp, "attempt to set config received 404 response")
		if err != nil {
			return err
		}

		result := model.WrappedResolvedAppConfig{}
		api.ParseResponse(resp, &result)
		printTable([]model.ResolvedAppConfig{*result.Record}, NewConfigRecord)
		return nil
	},
}

var configDeleteCmd = &cobra.Command{
	Use:          "delete KEY",
	Short:        "delete config",
	Long:         "Delete the config for the given app, env, stack, and key",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		logrus.Info(messageWithStackSuffix(
			fmt.Sprintf("deleting app config with key '%s' in environment '%s'", key, happyConfig.GetEnv()),
		))

		body := model.NewAppConfigLookupPayload(happyConfig.App(), happyConfig.GetEnv(), stack, key)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Delete(fmt.Sprintf("/v1/configs/%s", key), body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = inspectForErrors(resp, notFoundMessage)
		if err != nil {
			return err
		}

		result := model.WrappedResolvedAppConfig{}
		api.ParseResponse(resp, &result)
		if result.Record != nil {
			logrus.Infof("app config with key '%s' has been deleted", result.Record.Key)
		} else {
			logrus.Warn(notFoundMessage)
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
		fmt.Println("here in copy\n--------")

		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		srcAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("copying app config with key '%s' from %s to %s", key, srcAppEnvStack, destAppEnvStack)

		body := model.NewAppConfigDiffPayload(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Get("/v1/config/diff", body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		notFoundMessage := messageWithStackSuffix(
			fmt.Sprintf("app config with key '%s' could not be found in environment '%s'", key, happyConfig.GetEnv()),
		)
		err = inspectForErrors(resp, notFoundMessage)
		if err != nil {
			return err
		}

		result := model.WrappedResolvedAppConfig{}
		api.ParseResponse(resp, &result)
		if result.Record != nil {
			logrus.Infof("app config with key '%s' has been copied from %s to %s", result.Record.Key, srcAppEnvStack, destAppEnvStack)
		} else {
			logrus.Warn(notFoundMessage)
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
		happyConfig, err := api.GetHappyConfig(cmd)
		if err != nil {
			return err
		}

		srcAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), fromEnv, fromStack)
		destAppEnvStack := fmt.Sprintf("%s/%s/%s", happyConfig.App(), happyConfig.GetEnv(), stack)
		logrus.Infof("retrieving list of config keys that exist in %s and not %s", srcAppEnvStack, destAppEnvStack)

		body := model.NewAppConfigDiffPayload(happyConfig.App(), fromEnv, fromStack, happyConfig.GetEnv(), stack)
		client := client.NewHappyClient(happyConfig)
		resp, err := client.Get("/v1/config/diff", body)
		if err != nil {
			return errors.Wrap(err, "request failed with")
		}

		err = inspectForErrors(resp, "attempt to get config diff received 404 response")
		if err != nil {
			return err
		}

		result := model.ConfigDiffResponse{}
		api.ParseResponse(resp, &result)
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

func inspectForErrors(resp *http.Response, notFoundMessage string) error {
	if resp.StatusCode == http.StatusNotFound {
		return errors.New(notFoundMessage)
	} else if resp.StatusCode == http.StatusBadRequest {
		validationErrors := []api.ValidationError{}
		api.ParseResponse(resp, &validationErrors)
		message := ""
		for _, validationError := range validationErrors {
			message = message + fmt.Sprintf("\nhappy-api request failed with: %s", validationError.Message)
		}
		return errors.New(message)
	} else if resp.StatusCode != http.StatusOK {
		errorMessage := new(string)
		api.ParseResponse(resp, errorMessage)
		return errors.New(*errorMessage)
	}
	return nil
}
