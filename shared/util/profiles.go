package util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func GetAWSProfiles(options ...func(*config.LoadOptions)) ([]string, error) {
	profiles := []string{}
	loadOptions := config.LoadOptions{}
	for _, option := range options {
		option(&loadOptions)
	}

	configFile := config.DefaultSharedConfigFilename()
	if len(loadOptions.SharedConfigFiles) > 0 {
		configFile = loadOptions.SharedConfigFiles[0]
	}
	logrus.Infof("Loading profiles from %s", configFile)
	f, err := ini.Load(configFile)
	if err != nil {
		return profiles, errors.Wrapf(err, "unable to load %s", configFile)
	}

	for _, v := range f.Sections() {
		if strings.HasPrefix(v.Name(), "profile ") {
			profile, _ := strings.CutPrefix(v.Name(), "profile ")
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}
