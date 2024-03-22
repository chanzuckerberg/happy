package options

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/interfaces"
)

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
}

type DryRun string

const DryRunKey DryRun = "dry-run"

type EnableDynamoLocking string

const EnableDynamoLockingKey EnableDynamoLocking = "enable-dynamo-locking"

type EnableAppDebugLogsDuringDeployment struct{}

func DebugLogsDuringDeploymentFromCtx(ctx context.Context) bool {
	v, _ := ctx.Value(EnableAppDebugLogsDuringDeployment{}).(bool)
	return v
}

func NewDebugLogsDuringDeploymentCtx(ctx context.Context, enable bool) context.Context {
	return context.WithValue(ctx, EnableAppDebugLogsDuringDeployment{}, enable)
}
