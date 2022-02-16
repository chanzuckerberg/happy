package main

import (
	"time"

	"github.com/chanzuckerberg/happy/cmd"
	"github.com/chanzuckerberg/happy/pkg/log"
	"github.com/gen2brain/beeep"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&log.Formatter{})

	startTime := time.Now()
	if err := cmd.Execute(); err != nil {
		_ = beeep.Alert("Happy Error!", err.Error(), "assets/warning.png")
		logrus.Fatal(err)
	}
	if time.Since(startTime) > 30*time.Second {
		_ = beeep.Notify("Happy", "Successfully completed", "assets/information.png")
	}
}
