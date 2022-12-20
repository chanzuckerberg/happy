package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestTokenProvider struct{}

func (t TestTokenProvider) GetToken() (string, error) {
	return "test-token", nil
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
				authHeader, ok := req.Header["Authorization"]
				r.Equal(ok, true)
				r.Equal(authHeader, []string{"Bearer test-token"})

				contentHeader, ok := req.Header["Content-Type"]
				r.Equal(ok, true)
				r.Equal(contentHeader, []string{"application/json"})

				if tc.version != "undefined" {
					userAgentHeader, ok := req.Header["User-Agent"]
					r.Equal(ok, true)
					r.Equal(userAgentHeader, []string{fmt.Sprintf("%s/%s", clientName, tc.version)})
				}
			}))
			defer func() { testServer.Close() }()

			client := NewHappyClient(clientName, tc.version, testServer.URL, TestTokenProvider{})
			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			r.NoError(err)

			_, err = client.Do(req)
			r.NoError(err)
		})
	}
}
