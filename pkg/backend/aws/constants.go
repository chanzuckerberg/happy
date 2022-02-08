package aws

type TaskType string

const (
	TaskTypeDelete  TaskType = "delete"
	TaskTypeMigrate TaskType = "migrate"
)
