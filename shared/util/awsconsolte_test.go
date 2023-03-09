package util

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEmptyArn(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	_, err := Arn2ConsoleLink(linkOptions, "")
	r.Error(err)
}

func TestSecretsManagerLink(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	link, err := Arn2ConsoleLink(linkOptions, "arn:aws:secretsmanager:us-west-2:1234567890:secret:happy/test-secret-1A234B")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/secretsmanager/home?region=us-west-2#!/secret?name=happy%%2Ftest-secret", link)
}

func TestEcsClusterLink(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	link, err := Arn2ConsoleLink(linkOptions, "arn:aws:ecs:us-west-2:1234567890:cluster/happy-cluster")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/ecs/home?region=us-west-2#/clusters/happy-cluster/services", link)
}

func TestEcsServiceLink(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	link, err := Arn2ConsoleLink(linkOptions, "arn:aws:ecs:us-west-2:1234567890:service/happy-cluster/test-frontend")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/ecs/home?region=us-west-2#/clusters/happy-cluster/services/test-frontend/tasks", link)
}

func TestEcsTaskLink(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	link, err := Arn2ConsoleLink(linkOptions, "arn:aws:ecs:us-west-2:1234567890:task/happy-cluster/39d0c5743f794b46bcc2c5f3ebc1b5b0")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/ecs/home?region=us-west-2#/clusters/happy-cluster/tasks/39d0c5743f794b46bcc2c5f3ebc1b5b0/details", link)
}

func TestLogLinkFargate(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
		LaunchType:           LaunchTypeFargate,
	}
	link, err := Log2ConsoleLink(linkOptions, "group", "prefix", "containerid", "taskid")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logEventViewer:group=group;stream=prefix/containerid/taskid", link)
}

func TestLogLinkK8S(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
		LaunchType:           LaunchTypeK8S,
	}
	link, err := Log2ConsoleLink(linkOptions, "/rdev-eks/fluentbit-cloudwatch", "logStreamNameFilter=fluentbit-kube.var.log.containers", "my-app", "")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logsV2:log-groups/log-group/$252Frdev-eks$252Ffluentbit-cloudwatch$3FlogStreamNameFilter$3Dfluentbit-kube.var.log.containers.my-app", link)
}

func TestLogInsightsLinkK8S(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
		LaunchType:           LaunchTypeK8S,
	}
	queryId := uuid.New().String()
	expression := `fields @timestamp, log
| sort @timestamp desc
| limit 20
| filter kubernetes.namespace_name = "rdev-happy-env"
| filter kubernetes.pod_name like "myapp-frontend"`
	link, err := LogInsights2ConsoleLink(linkOptions, "/rdev-eks/fluentbit-cloudwatch", expression, queryId)
	r.NoError(err)

	desiredLink := fmt.Sprintf("https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logsV2:logs-insights$3FqueryDetail$3D~(end~0~start~-3600~timeType~'RELATIVE~unit~'seconds~editorString~'fields*20*40timestamp*2C*20log*0A*7C*20sort*20*40timestamp*20desc*0A*7C*20limit*2020*0A*7C*20filter*20kubernetes.namespace_name*20*3D*20*22rdev-happy-env*22*0A*7C*20filter*20kubernetes.pod_name*20like*20*22myapp-frontend*22~queryId~'%s~source~(~'*2Frdev-eks*2Ffluentbit-cloudwatch))", queryId)
	r.Equal(desiredLink, link)
}
