package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/blang/semver"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func MakeTestApp(r *require.Assertions) *APIApplication {
	cfg := setup.GetConfiguration()
	app := MakeApp(context.Background(), cfg)
	return app
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
					ver.Minor = ver.Minor + 1
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
			app := MakeTestApp(r)

			req := httptest.NewRequest("GET", "/versionCheck", nil)
			req.Header.Set(fiber.HeaderUserAgent, tc.userAgent)

			resp, err := app.FiberApp.Test(req)
			r.NoError(err)

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
					ver.Minor = ver.Minor - 1
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
			errorMessage: "expected version so be specified for happy-cli in the User-Agent header (format: happy-cli/<version>)",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(r)

			req := httptest.NewRequest("GET", "/versionCheck", nil)
			req.Header.Set(fiber.HeaderUserAgent, tc.userAgent)

			resp, err := app.FiberApp.Test(req)
			r.NoError(err)

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
