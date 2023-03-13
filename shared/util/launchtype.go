package util

import "strings"

type LaunchType string

const (
	LaunchTypeFargate LaunchType = "FARGATE"
	LaunchTypeEC2     LaunchType = "EC2"
	LaunchTypeK8S     LaunchType = "K8S"
)

func (l LaunchType) String() string {
	return strings.ToUpper(string(l))
}
