package options

import "github.com/chanzuckerberg/happy/pkg/interfaces"

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
	DryRun       bool
}
