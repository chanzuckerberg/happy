package setup

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if getAppEnv() == "test" {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

type OIDCProvider struct {
	IssuerURL string `mapstructure:"issuer_url"`
	ClientID  string `mapstructure:"client_id"`
}

type OIDCProviders []OIDCProvider

// example of how to parse multiple OIDC providers with one env var:
// "HAPI_AUTH_OIDC_PROVIDERS": "issuer1|clientid1,issuer2|clientid2"
func (p *OIDCProviders) UnmarshalText(text []byte) error {
	s := strings.Split(string(text), ",")
	for _, v := range s {
		m := strings.SplitN(v, "|", 2)
		if len(m) != 2 {
			return errors.Errorf("bad format of OIDCProviders %s", string(text))
		}

		*p = append(*p, OIDCProvider{
			IssuerURL: m[0],
			ClientID:  m[1],
		})
	}

	return nil
}

type AuthConfiguration struct {
	Enable    *bool         `mapstructure:"enable"`
	Providers OIDCProviders `mapstructure:"oidc_providers"`
}

type ApiConfiguration struct {
	Port     uint   `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
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
	Auth     AuthConfiguration     `mapstructure:"auth"`
	Api      ApiConfiguration      `mapstructure:"api"`
	Database DatabaseConfiguration `mapstructure:"database"`
}

func evaluateConfigWithEnvToTmp(configPath string) (string, error) {
	tmp, err := os.CreateTemp("./", "*.yaml")
	if err != nil {
		return "", errors.Wrap(err, "unable to create a temp config file")
	}

	cfile, err := os.Open(configPath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open %s", configPath)
	}

	_, err = evaluateConfigWithEnv(cfile, tmp)
	if err != nil {
		return "", errors.Wrap(err, "unable to populate the environment")
	}

	return tmp.Name(), nil
}

func envToMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			continue
		}
		envMap[s[0]] = s[1]
	}
	return envMap
}

// evaluateConfigWithEnv reads a configuration reader and injects environment variables
// that exist as part of the configuration in the form a go template. For example
// {{.ENV_VAR1}} will be replace with the value of the environment variable ENV_VAR1.
// Optional support for writting the contents to other places is supported by providing
// other writers. By default, the evaluated configuartion is returned as a reader.
func evaluateConfigWithEnv(configFile io.Reader, writers ...io.Writer) (io.Reader, error) {
	envMap := envToMap()

	b, err := io.ReadAll(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the config file")
	}

	t := template.New("appConfigTemplate")
	tmpl, err := t.Parse(string(b))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse template from: \n%s", string(b))
	}

	populated := []byte{}
	buff := bytes.NewBuffer(populated)
	writers = append(writers, buff)
	err = tmpl.Execute(io.MultiWriter(writers...), envMap)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute template")
	}
	return buff, nil
}

const defaultConfigYamlDir = "./"

func GetConfiguration() (*Configuration, error) {
	configYamlDir := defaultConfigYamlDir
	if len(os.Getenv("CONFIG_YAML_DIRECTORY")) > 0 {
		configYamlDir = os.Getenv("CONFIG_YAML_DIRECTORY")
	}
	path, err := filepath.Abs(configYamlDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get absolute path of %s", configYamlDir)
	}

	vpr := viper.New()
	appConfigFile := filepath.Join(path, "app-config.yaml")
	if _, err := os.Stat(appConfigFile); err == nil {
		tmp, err := evaluateConfigWithEnvToTmp(appConfigFile)
		if len(tmp) != 0 {
			defer os.Remove(tmp)
		}
		if err != nil {
			return nil, err
		}

		vpr.SetConfigFile(tmp)
		err = vpr.ReadInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config file")
		}
	}

	envConfigFilename := fmt.Sprintf("app-config.%s.yaml", getAppEnv())
	appEnvConfigFile := filepath.Join(path, envConfigFilename)
	if _, err := os.Stat(appEnvConfigFile); err == nil {
		tmp, err := evaluateConfigWithEnvToTmp(appEnvConfigFile)
		if len(tmp) != 0 {
			defer os.Remove(tmp)
		}
		if err != nil {
			return nil, err
		}

		vpr.SetConfigFile(tmp)
		err = vpr.MergeInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge env config")
		}
	}

	// unique case where we want to be able to configure static OIDC providers
	// but also provide them as an environment variable and have the two settings
	// combined (as in the two lists are combined). I suppose we could just give
	// priority to the environment variables and have them overwrite.
	envVpr := viper.New()
	envVpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	envVpr.SetEnvPrefix("HAPI")
	err = envVpr.BindEnv("auth.oidc_providers")
	if err != nil {
		return nil, errors.Wrap(err, "failed to bind environment auth.oidc_providers")
	}
	err = envVpr.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read in environment config")
	}
	envCfg := &Configuration{}
	err = envVpr.Unmarshal(envCfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration")
	}

	cfg := &Configuration{}
	err = vpr.Unmarshal(cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration")
	}
	cfg.Auth.Providers = append(cfg.Auth.Providers, envCfg.Auth.Providers...)
	// default to having auth enabled
	if cfg.Auth.Enable == nil {
		enable := true
		cfg.Auth.Enable = &enable
	}

	return cfg, nil
}

func getAppEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("DEPLOYMENT_STAGE")
	}
	return env
}
