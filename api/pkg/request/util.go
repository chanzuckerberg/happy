package request

import (
	"context"
	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type AWSCredentialsContextKey struct{}

func CtxWithAWSAuthHeaders(ctx *fiber.Ctx) (context.Context, error) {
	encodedAccessKeyId := ctx.Get("x-aws-access-key-id")
	encodedSecretAccessKey := ctx.Get("x-aws-secret-access-key")
	encodedSessionToken := ctx.Get("x-aws-session-token")
	return AddAWSAuthToCtx(ctx.Context(), encodedAccessKeyId, encodedSecretAccessKey, encodedSessionToken)
}

func AddAWSAuthToCtx(ctx context.Context, encodedAccessKeyId, encodedSecretAccessKey, encodedSessionToken string) (context.Context, error) {
	accessKeyId, err := b64.StdEncoding.DecodeString(encodedAccessKeyId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode base64 value from AWS access key id")
	}

	secretAccessKey, err := b64.StdEncoding.DecodeString(encodedSecretAccessKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode base64 value from AWS secret access key")
	}

	return context.WithValue(ctx, AWSCredentialsContextKey{}, AWSCredentials{
		AccessKeyID:     string(accessKeyId),
		SecretAccessKey: string(secretAccessKey),
		SessionToken:    encodedSessionToken, // session token should remain base64 encoded
	}), nil
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
