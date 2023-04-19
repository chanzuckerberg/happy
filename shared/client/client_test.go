package client

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

type TestTokenProvider struct{}

func (t TestTokenProvider) GetToken() (string, error) {
	return "test-token", nil
}

type AWSCredentialsProviderTest struct {
}

func (t AWSCredentialsProviderTest) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "pre-encoded-session-token",
	}, nil
}

func TestDoSuccess(t *testing.T) {
	testData := []struct {
		version string
	}{
		{
			version: "undefined",
		},
		{
			version: "1.0.1",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			clientName := "test"

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				authHeader := req.Header.Get("Authorization")
				r.Equal(authHeader, "Bearer test-token")

				accessKeyHeader := req.Header.Get("x-aws-access-key-id")
				r.Equal(accessKeyHeader, b64.StdEncoding.EncodeToString([]byte("access-key-id")))

				secretKeyHeader := req.Header.Get("x-aws-secret-access-key")
				r.Equal(secretKeyHeader, b64.StdEncoding.EncodeToString([]byte("secret-access-key")))

				sessionTokenHeader := req.Header.Get("x-aws-session-token")
				r.Equal(sessionTokenHeader, "pre-encoded-session-token")

				contentHeader := req.Header.Get("Content-Type")
				r.Equal(contentHeader, "application/json")

				if tc.version != "undefined" {
					userAgentHeader := req.Header.Get("User-Agent")
					r.Equal(userAgentHeader, fmt.Sprintf("%s/%s", clientName, tc.version))
				}
			}))
			defer func() { testServer.Close() }()

			client := NewHappyClient(clientName, tc.version, testServer.URL, TestTokenProvider{}, AWSCredentialsProviderTest{})
			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			r.NoError(err)

			_, err = client.Do(req)
			r.NoError(err)
		})
	}
}
