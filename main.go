package main

import (
	"github.com/chanzuckerberg/happy/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	log.SetLevel(log.InfoLevel)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
