package util

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/pkg/errors"
)

type LinkOptions struct {
	Region               string
	IntegrationSecretARN string
	LaunchType           LaunchType
	AWSAccountID         string
}

type LogReference struct {
	LinkOptions
	LogGroupName string
	Expression   string
}

func Arn2ConsoleLink(options LinkOptions, unparsedArn string) (string, error) {
	if len(unparsedArn) == 0 {
		return "", errors.New("ARN not provided")
	}
	if len(options.Region) == 0 {
		return "", errors.New("region not specified")
	}
	if len(options.IntegrationSecretARN) == 0 {
		return "", errors.New("integration secret ARN not specified")
	}
	resourceArn, err := arn.Parse(unparsedArn)
	if err != nil {
		return "", errors.Wrapf(err, "invalid ARN: %s", unparsedArn)
	}

	region := options.Region
	q := url.Values{
		"region": []string{region},
	}

	awsConsoleUrl := url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s.console.aws.amazon.com", region),
		RawQuery: q.Encode(),
	}

	switch resourceArn.Service {
	case "ecs":

		awsConsoleUrl.Path = "/ecs/home"

		resourceParts := strings.Split(resourceArn.Resource, "/")
		if len(resourceParts) < 2 {
			return "", errors.Errorf("ARN is not supported: %s", unparsedArn)
		}
		resourceType := resourceParts[0]
		resourceName := resourceParts[1]
		resourceSubName := ""
		if len(resourceParts) > 2 {
			resourceSubName = resourceParts[2]
		}

		switch resourceType {
		case "cluster":
			awsConsoleUrl.Fragment = fmt.Sprintf("/clusters/%s/services", resourceName)

			// Returns a link like this one:
			// fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services", region, region, resourceName), nil
			return awsConsoleUrl.String(), nil
		case "task":
			awsConsoleUrl.Fragment = fmt.Sprintf("/clusters/%s/tasks/%s/details", resourceName, resourceSubName)
			return awsConsoleUrl.String(), nil

			// Returns a link like this one:
			// fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/tasks/%s/details", region, region, resourceName, resourceSubName), nil
		case "service":
			awsConsoleUrl.Fragment = fmt.Sprintf("/clusters/%s/services/%s/tasks", resourceName, resourceSubName)
			return awsConsoleUrl.String(), nil

			// Returns a link like this one:
			// fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/tasks", region, region, resourceName, resourceSubName), nil
		}
		return "", errors.Errorf("resource %s is not supported", resourceType)

	case "secretsmanager":
		resourceParts := strings.Split(resourceArn.Resource, ":")
		resourceType := resourceParts[0]

		awsConsoleUrl.Path = "/secretsmanager/home"

		switch resourceType {
		case "secret":
			secretName := strings.ReplaceAll(url.QueryEscape(options.IntegrationSecretARN), "%", "%%")
			// HACK: .String() double encodes the ? in the fragment, resulting in %252F instead of %2F
			return awsConsoleUrl.String() + "#" + fmt.Sprintf("!/secret?name=%s", secretName), nil

			// Returns a link like this one:
			// fmt.Sprintf("https://%s.console.aws.amazon.com/secretsmanager/home?region=%s#!/secret?name=%s", region, region, secretName), nil
		}
		return "", errors.Errorf("resource %s is not supported", resourceType)
	}

	return "", errors.Errorf("service %s is not supported", unparsedArn)
}

func Log2ConsoleLink(options LinkOptions, logGroup string, logStreamPrefix string, containerName string, taskId string) (string, error) {
	if len(options.Region) == 0 {
		return "", errors.New("region not specified")
	}
	if len(logGroup) == 0 {
		return "", errors.New("logGroup not specified")
	}
	if len(logStreamPrefix) == 0 {
		return "", errors.New("logStreamPrefix not specified")
	}
	if len(containerName) == 0 {
		return "", errors.New("containerName not specified")
	}
	q := url.Values{
		"region": []string{options.Region},
	}

	if options.LaunchType == LaunchTypeK8S {
		encodedLogGroup := strings.ReplaceAll(url.QueryEscape(fmt.Sprintf("%s?%s.%s", url.QueryEscape(logGroup), logStreamPrefix, containerName)), "%", "$")
		awsConsoleUrl := url.URL{
			Scheme:   "https",
			Host:     fmt.Sprintf("%s.console.aws.amazon.com", options.Region),
			Path:     "/cloudwatch/home",
			RawQuery: q.Encode(),
			Fragment: fmt.Sprintf("logsV2:log-groups/log-group/%s", encodedLogGroup),
		}
		return awsConsoleUrl.String(), nil
	}

	if len(taskId) == 0 {
		return "", errors.New("taskId not specified")
	}

	awsConsoleUrl := url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s.console.aws.amazon.com", options.Region),
		Path:     "/cloudwatch/home",
		RawQuery: q.Encode(),
		Fragment: fmt.Sprintf("logEventViewer:group=%s;stream=%s/%s/%s", logGroup, logStreamPrefix, containerName, taskId),
	}

	// Returns a link like this one:
	// fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s#logEventViewer:group=%s;stream=%s/%s/%s", logRegion, logRegion, logGroup, logStreamPrefix, containerName, taskId)

	return awsConsoleUrl.String(), nil
}

func LogInsights2ConsoleLink(logReference LogReference, queryId string) (string, error) {
	if len(logReference.Region) == 0 {
		return "", errors.New("region not specified")
	}
	if len(logReference.Expression) == 0 {
		return "", errors.New("expression not specified")
	}

	q := url.Values{
		"region": []string{logReference.Region},
	}

	encodedQuery := strings.ReplaceAll(url.QueryEscape("?queryDetail="), "%", "$") + buildCloudWatchInsightsQuery(logReference.Expression, logReference.LogGroupName, queryId)
	awsConsoleUrl := url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s.console.aws.amazon.com", logReference.Region),
		Path:     "/cloudwatch/home",
		RawQuery: q.Encode(),
		Fragment: "logsV2:logs-insights",
	}
	return awsConsoleUrl.String() + encodedQuery, nil
}

func escape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(url.QueryEscape(s), "+", "%20"), "%", "*")
}

func buildCloudWatchInsightsQuery(expression string, logGroup string, queryId string) string {
	query := strings.Builder{}
	query.WriteString("~(end~0~start~-3600~timeType~'RELATIVE~unit~'seconds~editorString~'")
	query.WriteString(escape(expression))
	query.WriteString("~queryId~'")
	query.WriteString(queryId)
	query.WriteString("~source~(~'")
	query.WriteString(escape(logGroup))
	query.WriteString("))")

	return query.String()
}
