package options

import "github.com/chanzuckerberg/happy/pkg/cli/interfaces"

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
}
