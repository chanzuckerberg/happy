package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func sendGetRequest(app *APIApplication, route string, r *require.Assertions) *http.Response {
	svr := httptest.NewServer(app.mux)
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", svr.URL, route), nil)
	r.NoError(err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	r.NoError(err)

	return resp
}

func TestHealthSucceed(t *testing.T) {
	r := require.New(t)
	app := MakeTestApp(t)
	resp := sendGetRequest(app, "/v2/health", r)

	r.Equal(200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	r.NoError(err)

	jsonBody := map[string]interface{}{}
	err = json.Unmarshal(body, &jsonBody)
	r.NoError(err)

	r.Contains(jsonBody["status"], "ok")
}

func TestAppConfigsFail(t *testing.T) {
	testData := []struct {
		queryString   string
		errorResponse string
	}{
		{
			queryString:   "",
			errorResponse: "{\"code\":500,\"errors\":\"operation ListAppConfig: decode params: query: \\\"app_name\\\": field required\"}",
		},
		{
			queryString:   "environment=prod",
			errorResponse: "{\"code\":500,\"errors\":\"operation ListAppConfig: decode params: query: \\\"app_name\\\": field required\"}",
		},
		{
			queryString:   "app_name=testapp",
			errorResponse: "{\"code\":500,\"errors\":\"operation ListAppConfig: decode params: query: \\\"environment\\\": field required\"}",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)
			resp := sendGetRequest(app, fmt.Sprintf("/v2/app-configs?%s", tc.queryString), r)

			r.Equal(500, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			r.NoError(err)
			r.Equal(tc.errorResponse, body)
		})
	}
}
