package main

import (
	"context"
	"time"

	_ "github.com/chanzuckerberg/happy/api/docs" // import API docs
	"github.com/chanzuckerberg/happy/api/pkg/api"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	sentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func exec(ctx context.Context) error {
	cfg := setup.GetConfiguration()

	m, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	logrus.Info("Running with configuration:\n", string(m))

	err = sentry.Init(sentry.ClientOptions{
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

	return api.MakeFiberApp(ctx, cfg).Listen()
}

// @title       Happy API
// @description An API to encapsulate Happy Path functionality
// @BasePath    /
func main() {
	err := exec(context.Background())
	if err != nil {
		logrus.Error(err)
	}
}
