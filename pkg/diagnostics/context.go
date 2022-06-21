package diagnostics

import (
	"context"
	"time"

	"github.com/chanzuckerberg/happy/pkg/profiler"
	"github.com/pkg/errors"
)

type ContextKey string

const diagnosticsContextKey ContextKey = "diagnostics"
const warningsContextKey ContextKey = "warnings"
const profilerContextKey ContextKey = "performance profiling"
const tfeRunInfoKey ContextKey = "TFE run info"

var NotADiagnosticContextError = errors.New("not a diagnostic context")
var WarningsNotFoundError = errors.New("warnings not found")

type DiagnosticContext struct {
	context.Context
}

func ToDiagnosticContext(ctx context.Context) (DiagnosticContext, error) {
	if !isDiagnosticContext(ctx) {
		return DiagnosticContext{Context: ctx}, NotADiagnosticContextError
	}
	return DiagnosticContext{Context: ctx}, nil
}

func BuildDiagnosticContext(ctx context.Context) DiagnosticContext {
	ctx = context.WithValue(ctx, diagnosticsContextKey, "true")
	ctx = context.WithValue(ctx, warningsContextKey, &[]string{})
	ctx = context.WithValue(ctx, profilerContextKey, profiler.NewProfiler())
	ctx = context.WithValue(ctx, tfeRunInfoKey, NewTfeRunInfo())
	return DiagnosticContext{Context: ctx}
}

func isDiagnosticContext(ctx context.Context) bool {
	ok := ctx.Value(diagnosticsContextKey)
	return ok != nil && ok.(string) == "true"
}

func (dctx *DiagnosticContext) AddWarning(warning string) {
	warnings := dctx.Value(warningsContextKey).(*[]string)
	if warnings == nil {
		warnings = &[]string{}
	}
	*warnings = append(*warnings, warning)
}

func AddWarning(ctx context.Context, warning string) error {
	dctx, err := ToDiagnosticContext(ctx)
	if err != nil {
		return NotADiagnosticContextError
	}
	dctx.AddWarning(warning)
	return nil
}

func (dctx *DiagnosticContext) GetWarnings() ([]string, error) {
	warnings := dctx.Value(warningsContextKey)
	if warnings == nil {
		return []string{}, WarningsNotFoundError
	}
	return dedupeWarnings(*warnings.(*[]string)), nil
}

func GetWarnings(ctx context.Context) ([]string, error) {
	dctx, err := ToDiagnosticContext(ctx)
	if err != nil {
		return []string{}, NotADiagnosticContextError
	}
	return dctx.GetWarnings()
}

func dedupeWarnings(warnings []string) []string {
	keyMap := map[string]bool{}
	uniqueWarnings := []string{}
	for _, warning := range warnings {
		if _, ok := keyMap[warning]; !ok {
			uniqueWarnings = append(uniqueWarnings, warning)
			keyMap[warning] = true
		}
	}
	return uniqueWarnings
}

func getContextProfiler(ctx context.Context) (*profiler.Profiler, error) {
	contextProfiler := ctx.Value(profilerContextKey)
	if contextProfiler != nil {
		return contextProfiler.(*profiler.Profiler), nil
	} else {
		return nil, errors.New("Context does not have a profiler")
	}
}

func AddProfilerRuntime(ctx context.Context, startTime time.Time, sectorName string) {
	contextProfiler, err := getContextProfiler(ctx)
	if err == nil {
		contextProfiler.AddRuntime(startTime, sectorName)
	}
}

func PrintRuntimes(ctx context.Context) {
	contextProfiler, err := getContextProfiler(ctx)
	if err == nil {
		contextProfiler.PrintRuntimes()
	}
}
