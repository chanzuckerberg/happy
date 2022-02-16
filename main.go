package main

import (
	"github.com/chanzuckerberg/happy/cmd"
	"github.com/chanzuckerberg/happy/pkg/log"
	"github.com/gen2brain/beeep"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&log.Formatter{})

	if err := cmd.Execute(); err != nil {
		beeep.Alert("Happy Error!", err.Error(), "assets/warning.png")
		logrus.Fatal(err)
	}
}
