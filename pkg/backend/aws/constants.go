package aws

type TaskType string

const (
	TaskTypeDelete  TaskType = "delete"
	TaskTypeMigrate TaskType = "migrate"
)

type LaunchType string

const (
	LaunchTypeFargate LaunchType = "FARGATE"
	LaunchTypeECS     LaunchType = "EC2"
)
