package util

import (
	"context"
	"os"
)

type ctxKey int

const (
	keyIsCI ctxKey = iota
)

// BuildContext will populate context with basic information
// it is expected to be run from the root of the cli and thus
// values set here are available throughout the whole application
// NOTE: we need to wait till cobra 1.5.0
func BuildContext(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, keyIsCI, isCI())
	return ctx, nil
}

func isCI() bool {
	// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	return os.Getenv("CI") == "true"
}

func IsCI(ctx context.Context) bool {
	// isCI, ok := ctx.Value(keyIsCI).(bool)
	// return ok && isCI
	//HACK: until cobra 1.5.0 do the following instead
	return isCI()
}

type TaskStartTimeKey string

var CmdStartContextKey TaskStartTimeKey = "taskStartTime"
