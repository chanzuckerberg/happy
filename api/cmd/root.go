package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "happy-api",
		Short:         "happy-api - an API for Happy Path",
		Long:          `happy-api is the API backend that supports the Happy Path CLI.`,
		SilenceErrors: true,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
