package diagnostics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnosticContext(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()
	dctx := BuildDiagnosticContext(ctx, true)

	r.True(isDiagnosticContext(dctx))
	r.False(isDiagnosticContext(ctx))

	warnings, err := dctx.GetWarnings()
	r.NoError(err)
	r.Len(warnings, 0)

	dctx.AddWarning("test")

	warnings, err = dctx.GetWarnings()
	r.NoError(err)
	r.Len(warnings, 1)
	warnings, err = GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 1)
	_, err = GetWarnings(ctx)
	r.Error(err)

	_, err = GetWarnings(ctx)
	r.ErrorIs(err, NotADiagnosticContextError)
	r.Error(err)

	dctx1, err := ToDiagnosticContext(dctx)
	r.NoError(err)
	r.True(isDiagnosticContext(dctx1))
	dctx1.AddWarning("warning")

	warnings, err = GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 2)

	dctx1.AddWarning("warning")
	warnings, err = GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 2)

	err = AddWarning(dctx, "warning1")
	r.NoError(err)
	warnings, err = GetWarnings(dctx)
	r.NoError(err)
	r.Len(warnings, 3)

	_, err = ToDiagnosticContext(ctx)
	r.Error(err)

	r.True(IsInteractiveContext(dctx))
}
