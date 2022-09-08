package main

import (
	// import API docs
	_ "github.com/chanzuckerberg/happy-api/docs"
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

// @title       Happy API
// @description An API to encapsulate Happy Path functionality
// @host        localhost:3001
// @BasePath    /
func main() {
	err := exec()
	if err != nil {
		logrus.Error(err)
	}
}
