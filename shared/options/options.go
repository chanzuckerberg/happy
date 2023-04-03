package options

import "github.com/chanzuckerberg/happy/shared/interfaces"

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
}
