package util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func GetAwsProfiles() ([]string, error) {
	profiles := []string{}
	configFile := config.DefaultSharedConfigFilename()
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
