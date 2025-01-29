package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/blang/semver"
	"github.com/chanzuckerberg/happy/api/pkg/ent/enttest"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	mu sync.Mutex
)

func MakeTestApp(t *testing.T) *APIApplication {
	cfg := setup.GetConfiguration()

	// Even with a UUID in the data source name this is not thread safe so we need to use a mutex to prevent concurrent access
	mu.Lock()
	client := enttest.Open(t, "sqlite3", fmt.Sprintf("file:memdb%s?mode=memory&cache=shared&_fk=1", uuid.NewString()))
	mu.Unlock()

	testDB := store.MakeDB(cfg.Database).WithClient(client)
	app := MakeAPIApplication(context.Background(), cfg, testDB)
	return app
}

func sendVersionCheckRequest(r *require.Assertions, svr *httptest.Server, userAgent string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/versionCheck", svr.URL), nil)
	r.NoError(err)
	req.Header.Set(fiber.HeaderUserAgent, userAgent)

	client := http.DefaultClient
	resp, err := client.Do(req)
	r.NoError(err)

	return resp
}

func TestVersionCheckSucceed(t *testing.T) {
	testData := []struct {
		userAgent string
	}{
		{
			// with minimum version
			userAgent: fmt.Sprintf("happy-cli/%s", request.MinimumVersions["happy-cli"]),
		},
		{
			// with above minimum version
			userAgent: fmt.Sprintf(
				"happy-cli/%s",
				func() string {
					ver := semver.MustParse(request.MinimumVersions["happy-cli"])
					ver.Minor++
					return ver.String()
				}(),
			),
		},
		{
			// unrestricted client with a version
			userAgent: "other-client/0.0.0",
		},
		{
			// unrestricted client without a version
			userAgent: "other-client-without-version",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)
			svr := httptest.NewServer(app.mux)
			defer svr.Close()

			resp := sendVersionCheckRequest(r, svr, tc.userAgent)
			r.Equal(fiber.StatusOK, resp.StatusCode)
		})
	}
}

func TestVersionCheckFail(t *testing.T) {
	testData := []struct {
		userAgent    string
		errorMessage string
	}{
		{
			// with below minimum version
			userAgent: fmt.Sprintf(
				"happy-cli/%s",
				func() string {
					ver := semver.MustParse(request.MinimumVersions["happy-cli"])
					ver.Minor--
					return ver.String()
				}(),
			),
			errorMessage: "please upgrade your happy-cli client to at least",
		},
		{
			// restricted client with invalid version
			userAgent:    "happy-cli/foo",
			errorMessage: "No Major.Minor.Patch elements found",
		},
		{
			// unrestricted client without a version
			userAgent:    "happy-cli",
			errorMessage: "expected version to be specified for happy-cli in the User-Agent header (format: happy-cli/<version>)",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)
			svr := httptest.NewServer(app.mux)
			defer svr.Close()

			resp := sendVersionCheckRequest(r, svr, tc.userAgent)
			r.Equal(fiber.StatusBadRequest, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			r.NoError(err)

			jsonBody := map[string]interface{}{}
			err = json.Unmarshal(body, &jsonBody)
			r.NoError(err)

			r.Contains(jsonBody["message"], tc.errorMessage)
		})
	}
}
