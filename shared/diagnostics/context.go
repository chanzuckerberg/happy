package diagnostics

import (
	"context"
	"strconv"
	"time"

	"github.com/chanzuckerberg/happy/shared/profiler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ContextKey string

const diagnosticsContextKey ContextKey = "diagnostics"
const warningsContextKey ContextKey = "warnings"
const profilerContextKey ContextKey = "performance profiling"
const tfeRunInfoContextKey ContextKey = "TFE run info"
const interactiveContextKey ContextKey = "interactive"

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

func BuildDiagnosticContext(ctx context.Context, interactive bool) DiagnosticContext {
	ctx = context.WithValue(ctx, diagnosticsContextKey, "true")
	ctx = context.WithValue(ctx, warningsContextKey, &[]string{})
	ctx = context.WithValue(ctx, profilerContextKey, profiler.NewProfiler())
	ctx = context.WithValue(ctx, tfeRunInfoContextKey, NewTfeRunInfo())
	ctx = context.WithValue(ctx, interactiveContextKey, strconv.FormatBool(interactive))
	return DiagnosticContext{Context: ctx}
}

func isDiagnosticContext(ctx context.Context) bool {
	ok := ctx.Value(diagnosticsContextKey)
	return ok != nil && ok.(string) == "true"
}

func IsInteractiveContext(ctx context.Context) bool {
	ok := ctx.Value(interactiveContextKey)
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

func getContextTfeRunInfo(ctx context.Context) (*TfeRunInfo, error) {
	contextTfeRunInfo := ctx.Value(tfeRunInfoContextKey)
	if contextTfeRunInfo == nil {
		return nil, errors.New("Context does not have TFE run info")
	}
	return contextTfeRunInfo.(*TfeRunInfo), nil
}

func AddTfeRunInfoUrl(ctx context.Context, url string) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to add TFE url: %s", err.Error())
		return
	}
	info.AddTfeUrl(url)
}

func AddTfeRunInfoOrg(ctx context.Context, org string) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to add TFE org: %s", err.Error())
		return
	}
	info.AddOrg(org)
}

func AddTfeRunInfoWorkspace(ctx context.Context, workspace string) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to add TFE workspace: %s", err.Error())
		return
	}
	info.AddWorkspace(workspace)
}

func AddTfeRunInfoRunId(ctx context.Context, runId string) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to add TFE run ID: %s", err.Error())
		return
	}
	info.AddRunId(runId)
}

func PrintTfeRunLink(ctx context.Context) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to print TFE run link: %s", err.Error())
		return
	}
	info.PrintTfeRunLink()
}

func GetTfeRunLink(ctx context.Context) (string, error) {
	info, err := getContextTfeRunInfo(ctx)
	if err != nil {
		logrus.Debugf("Unable to print TFE run link: %s", err.Error())
		return "", err
	}
	return info.MakeTfeRunLink()
}

func getContextProfiler(ctx context.Context) (*profiler.Profiler, error) {
	contextProfiler := ctx.Value(profilerContextKey)
	if contextProfiler == nil {
		return nil, errors.New("Context does not have a profiler")
	}
	return contextProfiler.(*profiler.Profiler), nil
}

func AddProfilerRuntime(ctx context.Context, startTime time.Time, sectorName string) {
	contextProfiler, err := getContextProfiler(ctx)
	if err != nil {
		logrus.Debugf("Unable to add profiler runtime: %s", err.Error())
		return
	}
	contextProfiler.AddRuntime(startTime, sectorName)
}

func PrintRuntimes(ctx context.Context) {
	contextProfiler, err := getContextProfiler(ctx)
	if err != nil {
		logrus.Debugf("Unable to print profiler runtimes: %s", err.Error())
		return
	}
	contextProfiler.PrintRuntimes()
}
