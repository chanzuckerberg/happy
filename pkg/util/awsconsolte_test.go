package util

import (
	"testing"

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

func TestLogLink(t *testing.T) {
	r := require.New(t)

	linkOptions := LinkOptions{
		Region:               "us-west-2",
		IntegrationSecretARN: "happy/test-secret",
	}
	link, err := Log2ConsoleLink(linkOptions, "group", "prefix", "containerid", "taskid")
	r.NoError(err)
	r.Equal("https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logEventViewer:group=group;stream=prefix/containerid/taskid", link)
}
