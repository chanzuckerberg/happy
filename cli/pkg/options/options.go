package options

import "github.com/chanzuckerberg/happy/cli/pkg/interfaces"

type WaitOptions struct {
	StackName    string
	Orchestrator interfaces.OrchestratorInterface
	Services     []string
}
