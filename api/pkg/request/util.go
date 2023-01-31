package request

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type AWSCredentialsContextKey struct{}

func CtxWithAWSAuthHeaders(ctx *fiber.Ctx) context.Context {
	return context.WithValue(ctx.Context(), AWSCredentialsContextKey{}, AWSCredentials{
		AccessKeyID:     ctx.Get("x-aws-access-key-id"),
		SecretAccessKey: ctx.Get("x-aws-secret-access-key"),
		SessionToken:    ctx.Get("x-aws-session-token"),
	})
}
