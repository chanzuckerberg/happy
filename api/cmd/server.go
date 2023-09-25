package cmd

import (
	"context"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/api"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "run the happy api server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec(context.Background())
	},
}

func exec(ctx context.Context) error {
	cfg := setup.GetConfiguration()
	cfg.LogConfiguration()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.DSN,
		Environment:      cfg.Api.DeploymentStage,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err == nil {
		logrus.Info("Sentry enabled for environment: ", cfg.Api.DeploymentStage)
		// Flush buffered events before the program terminates.
		// Set the timeout to the maximum duration the program can afford to wait.
		defer sentry.Flush(2 * time.Second)
	} else {
		logrus.Info("Sentry disabled for environment: ", cfg.Api.DeploymentStage)
	}

	return api.MakeApp(ctx, cfg).Listen()
}
