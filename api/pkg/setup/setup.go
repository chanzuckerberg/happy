package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if os.Getenv("APP_ENV") == "test" {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

type ApiConfiguration struct {
	Port      uint   `mapstructure:"port"`
	LogLevel  string `mapstructure:"log_level"`
	IssuerURL string `mapstructure:"oidc_issuer_url"`
	ClientID  string `mapstructure:"oidc_client_id"`
}

type DBDriver string

const (
	Sqlite   DBDriver = "sqlite"
	Postgres DBDriver = "postgres"
)

type DatabaseConfiguration struct {
	Driver         DBDriver `mapstructure:"driver"`
	DataSourceName string   `mapstructure:"data_source_name"`
	LogLevel       string   `mapstructure:"log_level"`
}

type Configuration struct {
	Api      ApiConfiguration      `mapstructure:"api"`
	Database DatabaseConfiguration `mapstructure:"database"`
}

func GetConfiguration() (*Configuration, error) {
	configYamlDir := os.Getenv("CONFIG_YAML_DIRECTORY")
	path, err := filepath.Abs(configYamlDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get absolute path of %s", configYamlDir)
	}

	vpr := viper.New()
	vpr.SetEnvPrefix("happy_api")
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindVars := []string{"api.oidc_issuer_url", "api.oidc_client_id"}
	for _, bindVar := range bindVars {
		err = vpr.BindEnv(bindVar)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to bind to env var with '%s'", bindVar)
		}
	}

	appConfigFile := filepath.Join(path, "app-config.yaml")
	if _, err := os.Stat(appConfigFile); err == nil {
		vpr.SetConfigFile(appConfigFile)
		err = vpr.ReadInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config file")
		}
	}

	envConfigFilename := fmt.Sprintf("app-config.%s.yaml", os.Getenv("APP_ENV"))
	appEnvConfigFile := filepath.Join(path, envConfigFilename)
	if _, err := os.Stat(appEnvConfigFile); err == nil {
		vpr.SetConfigFile(appEnvConfigFile)
		err = vpr.MergeInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge env config")
		}
	}

	cfg := &Configuration{}
	err = vpr.Unmarshal(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration")
	}

	return cfg, nil
}
