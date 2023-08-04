package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

// Unset an environment variable by key
func (e *environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

// IsSet returns whether or not a key is currently set in the environ
func (e *environ) IsSet(key string) bool {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			return true
		}
	}
	return false
}

// Set adds an environment variable, replacing any existing ones of the same key
func (e *environ) Set(key, val string) {
	e.Unset(key)
	*e = append(*e, key+"="+val)
}

var configExecCmd = &cobra.Command{
	Use:          "exec",
	Short:        "execute a comment with configs",
	Long:         "Execute a command with configs for the given app, env, and stack",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		dashIx := cmd.ArgsLenAtDash()
		if dashIx == -1 {
			return errors.New("please separate services and command with '--'. See usage")
		}
		if err := cobra.MinimumNArgs(1)(cmd, args[dashIx:]); err != nil {
			return errors.Wrap(err, "must specify command to run. See usage")
		}
		return nil
	},
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
			fmt.Sprintf("executing with app configs in environment '%s'", happyClient.HappyConfig.GetEnv()),
		))

		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
		result, err := api.ListConfigs(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stack)
		if err != nil {
			if errors.Is(err, client.ErrRecordNotFound) {
				return errors.New("attempt to fetch configs received 404 response")
			}
			return err
		}

		dashIx := cmd.ArgsLenAtDash()
		command, commandArgs := args[dashIx], args[dashIx+1:]

		env := environ(os.Environ())
		for _, rawSecret := range result.Records {
			envVarKey := rawSecret.Key
			if env.IsSet(envVarKey) {
				fmt.Fprintf(os.Stderr, "warning: overwriting environment variable %s\n", envVarKey)
			}
			env.Set(envVarKey, rawSecret.Value)

		}
		return exec(command, commandArgs, env)
	},
}
