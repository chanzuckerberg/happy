package options

import "github.com/chanzuckerberg/happy/shared/interfaces"

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
}

type DryRun string

const DryRunKey DryRun = "dry-run"

type EnableDynamoLocking string

const EnableDynamoLockingKey EnableDynamoLocking = "enable-dynamo-locking"
