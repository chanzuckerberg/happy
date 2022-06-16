package stack_mgr

import (
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type StackManagementOptions struct {
	StackMeta    *StackMeta
	Stack        *Stack
	StackService *StackService
	HappyConfig  *config.HappyConfig
	StackTags    map[string]string
	Backend      *backend.Backend
	StackName    string
	DryRun       bool
}

func NewStackManagementOptions(stackName string) *StackManagementOptions {
	return &StackManagementOptions{StackName: stackName, StackTags: map[string]string{}}
}

func (o *StackManagementOptions) WithStackMeta(stackMeta *StackMeta) *StackManagementOptions {
	o.StackMeta = stackMeta
	return o
}

func (o *StackManagementOptions) WithStack(stack *Stack) *StackManagementOptions {
	o.Stack = stack
	return o
}

func (o *StackManagementOptions) WithHappyConfig(happyConfig *config.HappyConfig) *StackManagementOptions {
	o.HappyConfig = happyConfig
	return o
}

func (o *StackManagementOptions) WithStackTags(stackTags map[string]string) *StackManagementOptions {
	o.StackTags = stackTags
	return o
}

func (o *StackManagementOptions) WithStackService(stackService *StackService) *StackManagementOptions {
	o.StackService = stackService
	return o
}

func (o *StackManagementOptions) WithBackend(backend *backend.Backend) *StackManagementOptions {
	o.Backend = backend
	return o
}

func (o *StackManagementOptions) WithDryRun(dryRun bool) *StackManagementOptions {
	o.DryRun = dryRun
	return o
}
