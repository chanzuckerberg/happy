package main

import (
	"github.com/chanzuckerberg/happy/api/cmd"
	_ "github.com/chanzuckerberg/happy/api/docs" // import API docs
	"github.com/sirupsen/logrus"
)

// @title       Happy API
// @description An API to encapsulate Happy Path functionality
// @BasePath    /
func main() {
	err := cmd.Execute()
	if err != nil {
		logrus.Error(err)
	}
}
