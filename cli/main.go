package main

import (
	"time"

	"github.com/chanzuckerberg/happy/cli/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/log"
	"github.com/chanzuckerberg/happy/cli/pkg/output"
	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		logrus.Debugf("Successfully loaded environment variables from .env")
	} else {
		logrus.Debugf("Did not load environment variable files .env (%s), moving on", err.Error())
	}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&log.Formatter{})

	startTime := time.Now()
	if err = cmd.Execute(); err != nil {
		if cmd.Interactive {
			_ = beeep.Alert("Happy Error!", err.Error(), "assets/warning.png")
		}

		printer := output.NewPrinter(cmd.OutputFormat)
		printer.Fatal(err)

		return
	}

	if cmd.Interactive && time.Since(startTime) > 30*time.Second {
		_ = beeep.Notify("Happy", "Successfully completed", "assets/information.png")
	}
}
