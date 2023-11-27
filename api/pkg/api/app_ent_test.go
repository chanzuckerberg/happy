package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chanzuckerberg/happy/api/pkg/ent/enttest"
	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func MakeTestOgentApp(t *testing.T) *ogent.Server {
	cfg := setup.GetConfiguration()

	// Even with a UUID in the data source name this is not thread safe so we need to use a mutex to prevent concurrent access
	mu.Lock()
	client := enttest.Open(t, "sqlite3", fmt.Sprintf("file:memdb%s?mode=memory&cache=shared&_fk=1", uuid.NewString()))
	mu.Unlock()

	testDB := store.MakeDB(cfg.Database).WithClient(client)
	app, err := MakeOgentServerWithDB(context.Background(), cfg, testDB)
	r := require.New(t)
	r.NoError(err)
	return app
}

func TestHealthSucceed(t *testing.T) {
	r := require.New(t)
	app := MakeTestOgentApp(t)

	req := httptest.NewRequest(http.MethodGet, "/v2/health", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	res := w.Result()

	r.Equal(200, res.StatusCode)
	r.Equal("{\"status\":\"ok\"}", w.Body.String())
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
			app := MakeTestOgentApp(t)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v2/app-configs?%s", tc.queryString), nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			res := w.Result()

			r.Equal(500, res.StatusCode)
			r.Equal(tc.errorResponse, w.Body.String())
		})
	}
}
