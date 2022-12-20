package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestTokenProvider struct{}

func (t TestTokenProvider) GetToken() (string, error) {
	return "test-token", nil
}

func TestAddAuthSuccess(t *testing.T) {
	r := require.New(t)
	client := NewHappyClient("test", "0.0.0", "https://fake.hapi.io", TestTokenProvider{})

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", client.apiBaseUrl, "/test"), nil)
	r.NoError(err)

	err = client.addAuth(req)
	r.NoError(err)

	authHeader, ok := req.Header["Authorization"]

	r.Equal(ok, true)
	r.Equal(authHeader, []string{"Bearer test-token"})
}
