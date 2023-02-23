package util

type ManagedResource struct {
	Name      string
	Type      string
	Provider  string
	Module    string
	Instances []string
	ManagedBy string
}
