package diagnostics

import (
	"context"

	"github.com/pkg/errors"
)

const diagnosticsContextKey = "diagnostics"
const warningsContextKey = "warnings"

type DiagnosticContext struct {
	context.Context
}

func ToDiagnosticContext(ctx context.Context) (DiagnosticContext, error) {
	if !IsDiagnosticContext(ctx) {
		return DiagnosticContext{Context: ctx}, errors.New("not a diagnostic context")
	}
	return DiagnosticContext{Context: ctx}, nil
}

func BuildDiagnosticContext(ctx context.Context) DiagnosticContext {
	ctx = context.WithValue(ctx, diagnosticsContextKey, "true")
	ctx = context.WithValue(ctx, warningsContextKey, &[]string{})
	return DiagnosticContext{Context: ctx}
}

func IsDiagnosticContext(ctx context.Context) bool {
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

func (dctx *DiagnosticContext) GetWarnings() []string {
	return *dctx.Value(warningsContextKey).(*[]string)
}

func GetWarnings(ctx context.Context) ([]string, error) {
	if !IsDiagnosticContext(ctx) {
		return []string{}, errors.New("not a diagnostic context")
	}
	warnings := ctx.Value(warningsContextKey)
	if warnings == nil {
		return []string{}, errors.New("warnings not found")
	}
	return *warnings.(*[]string), nil
}
