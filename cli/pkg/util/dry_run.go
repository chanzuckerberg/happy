package util

type DryRunType bool

const (
	CompleteRun = DryRunType(true)
	DryRun      = DryRunType(false)
)
