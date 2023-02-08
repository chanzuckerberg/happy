package request

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
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

type AWSCredentialsProvider struct{}

func (c AWSCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	val := ctx.Value(AWSCredentialsContextKey{}).(AWSCredentials)
	return aws.Credentials{
		AccessKeyID:     val.AccessKeyID,
		SecretAccessKey: val.SecretAccessKey,
		SessionToken:    val.SessionToken,
	}, nil
}

func MakeCredentialProvider(ctx context.Context) aws.CredentialsProvider {
	return AWSCredentialsProvider{}
}
