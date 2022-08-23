package main

import (
	"github.com/chanzuckerberg/happy-api/pkg/api"
	"github.com/sirupsen/logrus"
)

func exec() error {
	app, err := api.MakeApp()
	if err != nil {
		return err
	}
	return app.Listen(":3001")
}

func main() {
	err := exec()
	if err != nil {
		logrus.Error(err)
	}
}
