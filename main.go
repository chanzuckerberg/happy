package main

import (
	"time"

	"github.com/chanzuckerberg/happy/cmd"
	"github.com/chanzuckerberg/happy/pkg/log"
	"github.com/chanzuckerberg/happy/pkg/output"
	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Infof("failed to load environment variables from .env: %s", err.Error())
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
