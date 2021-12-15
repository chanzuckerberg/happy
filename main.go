package main

import (
	"github.com/chanzuckerberg/happy-deploy/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
