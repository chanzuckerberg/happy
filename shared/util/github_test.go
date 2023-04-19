package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGithubDeploymentsQuery(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/json")

		fileName := fmt.Sprintf("./testdata%s.%s.json", r.URL.String(), r.Method)

		f, err := os.Open(fileName)
		req.NoError(err)
		_, err = io.Copy(w, f)
		req.NoError(err)

		w.WriteHeader(204)
	}))
	defer ts.Close()

	sha, err := GetLatestSuccessfulDeployment(context.Background(), ts.URL+"/graphql", "token", "staging", "owner", "repo")
	req.NoError(err)
	req.Equal("bd7b4339", sha)
}
