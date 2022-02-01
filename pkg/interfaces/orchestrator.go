package interfaces

type OrchestratorInterface interface {
	GetEvents(stack string, services []string) error
}
