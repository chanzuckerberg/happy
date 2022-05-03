package diagnostics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnosticContext(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()
	dctx := BuildDiagnosticContext(ctx)

	r.True(IsDiagnosticContext(dctx))
	r.False(IsDiagnosticContext(ctx))

	dctx.AddWarning("test")

	warnings := dctx.GetWarnings()
	r.Len(warnings, 1)
	warnings, err := GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 1)
	_, err = GetWarnings(ctx)
	r.Error(err)

	dctx1, err := ToDiagnosticContext(dctx)
	r.NoError(err)
	r.True(IsDiagnosticContext(dctx1))
	dctx1.AddWarning("warning")

	warnings, err = GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 2)

	_, err = ToDiagnosticContext(ctx)
	r.Error(err)
}
