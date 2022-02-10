package config

import "strings"

type LaunchType string

const (
	LaunchTypeFargate LaunchType = "FARGATE"
	LaunchTypeEC2     LaunchType = "EC2"
)

func (l LaunchType) String() string {
	return strings.ToUpper(string(l))
}
