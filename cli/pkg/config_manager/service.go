package config_manager

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
)

func configureService(bootstrapConfig *config.Bootstrap, dockerPath, defaultServicePort string) (*Service, error) {
	dockerFileName := filepath.Base(dockerPath)
	contextPath, err := filepath.Rel(bootstrapConfig.HappyProjectRoot, filepath.Dir(dockerPath))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain relative path")
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Would you like to use dockerfile %s/%s as a service in your stack?", contextPath, dockerFileName),
	}
	err = survey.AskOne(prompt, &confirm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prompt")
	}
	if !confirm {
		return nil, ErrServiceSkipped
	}

	serviceName := ""
	prompt1 := &survey.Input{
		Message: fmt.Sprintf("What would you like to name the service for %s/%s?", contextPath, dockerFileName),
		Help:    "This will be the name of the service in your stack, lowercased and hyphenated",
		Default: normalizeKey(filepath.Base(filepath.Dir(dockerPath))),
	}
	err = survey.AskOne(prompt1, &serviceName,
		survey.WithValidator(survey.Required),
		survey.WithValidator(survey.MinLength(3)),
		survey.WithValidator(survey.MaxLength(15)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to prompt")
	}

	serviceName = normalizeKey(serviceName)

	serviceType := serviceTypePrivate
	prompt2 := &survey.Select{
		Message: fmt.Sprintf("What kind of service is %s?", serviceName),
		Options: []string{
			serviceTypeExternal,
			serviceTypeInternal,
			serviceTypePrivate,
		},
		Default: serviceType,
	}

	err = survey.AskOne(prompt2, &serviceType)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain an aws profile")
	}

	var ok bool
	serviceType, ok = serviceTypeMapping[serviceType]
	if !ok {
		return nil, errors.Wrapf(err, "unexpected service type")
	}

	port := ""
	prompt3 := &survey.Input{
		Message: fmt.Sprintf("Which port does service %s listen on?", serviceName),
		Default: defaultServicePort,
	}
	err = survey.AskOne(prompt3, &port,
		survey.WithValidator(survey.Required),
		survey.WithValidator(survey.MinLength(2)),
		survey.WithValidator(survey.MaxLength(5)),
		survey.WithValidator(PortValidator))
	if err != nil {
		return nil, errors.Wrap(err, "unable to prompt")
	}

	uri := ""
	prompt3 = &survey.Input{
		Message: fmt.Sprintf("Which uri does %s respond on?", serviceName),
		Help:    "This is the relative path that the service will respond on, e.g. /api/v1",
		Default: "/",
	}
	err = survey.AskOne(prompt3, &uri,
		survey.WithValidator(survey.Required),
		survey.WithValidator(survey.MinLength(1)),
		survey.WithValidator(survey.MaxLength(255)),
		survey.WithValidator(URIValidator))
	if err != nil {
		return nil, errors.Wrap(err, "unable to prompt")
	}

	uri, _ = strings.CutSuffix(uri, "/")
	portNumber, err := strconv.Atoi(port)

	if err != nil {
		return nil, errors.Wrap(err, "port number is not valid")
	}

	return &Service{
		Name:           serviceName,
		ServiceType:    serviceType,
		DockerfilePath: dockerFileName,
		Context:        contextPath,
		Port:           portNumber,
		Uri:            uri,
	}, nil
}
