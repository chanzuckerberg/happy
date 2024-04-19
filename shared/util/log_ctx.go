package util

import "context"

type logGroupCtx struct{}

func NewLogGroupContext(ctx context.Context, logGroupPrefix string) context.Context {
	return context.WithValue(ctx, logGroupCtx{}, logGroupPrefix)
}

func LogGroupFromContext(ctx context.Context) string {
	logGroupPrefix, _ := ctx.Value(logGroupCtx{}).(string)
	return logGroupPrefix
}
