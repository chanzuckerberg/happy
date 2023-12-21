package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

// copied from request/auth_test.go
// I guess there are weird compile issues when sharing functions/structures across test files
func newDummyJWT(r *require.Assertions, subject, email string) string {
	mapClaims := jwt.MapClaims{
		"sub":   subject,
		"email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	ss, err := token.SignedString([]byte{})
	r.NoError(err)
	return ss
}

type reqData struct {
	body        map[string]interface{}
	queryParams map[string]string
	headers     map[string]string
}

func createRequest(svr *httptest.Server, method, route string, data reqData, r *require.Assertions) *http.Request {
	reader := &bytes.Reader{}
	values := url.Values{}
	for k, v := range data.queryParams {
		values.Add(k, v)
	}
	query := values.Encode()
	route = fmt.Sprintf("%s?%s", route, query)

	body, err := json.Marshal(data.body)
	r.NoError(err)
	reader = bytes.NewReader(body)

	// if method == http.MethodGet {
	// 	queryBytes, err := urlquery.Marshal(bodyMap)
	// 	r.NoError(err)
	// 	queryString = string(queryBytes)
	// 	route = fmt.Sprintf("%s?%s", route, queryString)
	// } else {
	// 	body, err := json.Marshal(bodyMap)
	// 	r.NoError(err)
	// 	reader = bytes.NewReader(body)
	// }
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", svr.URL, route), reader)
	r.NoError(err)

	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", newDummyJWT(r, "subject", "email@email.com")))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	for k, v := range data.headers {
		req.Header.Set(k, v)
	}
	return req
}
func makeRequest(svr *httptest.Server, method, route string, data reqData, r *require.Assertions) *http.Response {
	req := createRequest(svr, method, route, data, r)
	client := http.DefaultClient
	resp, err := client.Do(req)
	r.NoError(err)
	return resp
}

func makeSuccessfulRequest(app *APIApplication, method, route string, data reqData, r *require.Assertions) map[string]interface{} {
	svr := httptest.NewServer(app.mux)
	defer svr.Close()

	resp := makeRequest(svr, method, route, data, r)
	// r.Equal(fiber.StatusOK, resp.StatusCode)
	fmt.Println("...resp.StatusCode", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	r.NoError(err)
	fmt.Println("...resp.Body", string(body))

	jsonBody := map[string]interface{}{}
	err = json.Unmarshal(body, &jsonBody)
	r.NoError(err)

	return jsonBody
}

func makeInvalidRequest(app *APIApplication, method, route string, data reqData, r *require.Assertions) []map[string]interface{} {
	svr := httptest.NewServer(app.mux)
	defer svr.Close()

	resp := makeRequest(svr, method, route, data, r)
	r.Equal(fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	r.NoError(err)

	jsonBody := []map[string]interface{}{}
	err = json.Unmarshal(body, &jsonBody)
	r.NoError(err)

	return jsonBody
}

func TestSetConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		reqData      reqData
		expectRecord map[string]interface{}
	}{
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "bar",
					"key":         "TEST",
					"value":       "test-val",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
			},
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"key":         "TEST",
					"value":       "test-val2",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "rdev",
				"key":         "TEST",
				"value":       "test-val2",
			},
		},
		{
			// test that special characters are standardized
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"key":         "TEST-2*()$",
					"value":       "test-val2",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "rdev",
				"key":         "TEST_2____",
				"value":       "test-val2",
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			respBody := makeSuccessfulRequest(app, http.MethodPost, "/v1/configs", tc.reqData, r)

			record := respBody["record"].(map[string]interface{})

			_, createdAtPresent := record["created_at"]
			r.Equal(true, createdAtPresent)
			_, updatedAtPresent := record["updated_at"]
			r.Equal(true, updatedAtPresent)

			for _, key := range []string{"id", "created_at", "updated_at"} {
				r.NotNil(record[key])
				delete(record, key)
			}

			r.EqualValues(tc.expectRecord, record)
		})
	}
}

func TestSetConfigRouteFailure(t *testing.T) {
	testData := []struct {
		reqData     reqData
		failedField string
	}{
		{
			reqData: reqData{
				body: map[string]interface{}{
					"environment": "rdev",
					"stack":       "bar",
					"key":         "TEST",
					"value":       "test-val",
				},
			},
			failedField: "app_name",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name": "testapp",
					"stack":    "bar",
					"key":      "TEST",
					"value":    "test-val",
				},
			},
			failedField: "environment",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "bar",
					"value":       "test-val",
				},
			},
			failedField: "key",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "bar",
					"key":         "TEST",
				},
			},
			failedField: "value",
		},
		{
			// with invalid environment value
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "something",
					"stack":       "bar",
					"key":         "TEST",
				},
			},
			failedField: "environment",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			respBody := makeInvalidRequest(app, http.MethodPost, "/v1/configs", tc.reqData, r)

			r.Equal(tc.failedField, respBody[0]["failed_field"])
		})
	}
}

func TestSetConfigRouteFailsWithMalformedValue(t *testing.T) {
	testData := []struct {
		reqData      reqData
		errorMessage string
	}{
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    13,
					"environment": "rdev",
					"stack":       "bar",
					"key":         "TEST",
					"value":       "test-val",
				},
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.app_name of type string",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": 13,
					"stack":       "bar",
					"key":         "TEST",
					"value":       "test-val",
				},
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.environment of type string",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       13,
					"key":         "TEST",
					"value":       "test-val",
				},
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.stack of type string",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "",
					"key":         13,
					"value":       "test-val",
				},
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.key of type string",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "",
					"key":         "TEST",
					"value":       13,
				},
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.value of type string",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			respBody := makeInvalidRequest(app, http.MethodPost, "/v1/configs", tc.reqData, r)

			r.Contains(respBody[0]["message"], tc.errorMessage)
		})
	}
}

func TestGetConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds        []*model.AppConfigPayload
		reqData      reqData
		expectRecord map[string]interface{}
	}{
		{
			// only env config exists, looking up by env -> returns env config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "dev", "", "TEST", "test-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "dev",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "dev",
				"key":         "TEST",
				"value":       "test-val",
				"source":      "environment",
			},
		},
		{
			// only env config exists, looking up by stack -> returns env config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "dev", "", "TEST", "test-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "dev",
					"stack":       "bar",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "dev",
				"key":         "TEST",
				"value":       "test-val",
				"source":      "environment",
			},
		},
		{
			// env and stack configs exists, looking up by stack -> returns stack config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "dev", "", "TEST", "test-val"),
				model.NewAppConfigPayload("testapp", "dev", "bar", "TEST", "test-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "dev",
					"stack":       "bar",
				},
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "dev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
				"source":      "stack",
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			for _, input := range tc.seeds {
				_, err := cmd.MakeConfig(app.DB).SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, http.MethodGet, "/v1/configs/TEST", tc.reqData, r)
			record := respBody["record"].(map[string]interface{})

			_, createdAtPresent := record["created_at"]
			r.Equal(true, createdAtPresent)
			_, updatedAtPresent := record["updated_at"]
			r.Equal(true, updatedAtPresent)

			for _, key := range []string{"id", "created_at", "updated_at"} {
				r.NotNil(record[key])
				delete(record, key)
			}

			r.EqualValues(tc.expectRecord, record)
		})
	}
}

func TestDeleteConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []*model.AppConfigPayload
		reqData       reqData
		expectRecord  map[string]interface{}
		expectDeleted bool
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
				},
			},
			expectRecord:  nil,
			expectDeleted: false,
		},
		{
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "TEST", "test-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
				},
			},
			expectRecord: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "foo",
				"key":         "TEST",
				"value":       "test-val",
			},
			expectDeleted: true,
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			for _, input := range tc.seeds {
				_, err := cmd.MakeConfig(app.DB).SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, http.MethodDelete, "/v1/configs/TEST", tc.reqData, r)

			if tc.expectRecord == nil {
				r.Nil(respBody["record"])
			} else {
				record := respBody["record"].(map[string]interface{})

				_, createdAtPresent := record["created_at"]
				r.Equal(true, createdAtPresent)
				_, updatedAtPresent := record["updated_at"]
				r.Equal(true, updatedAtPresent)

				delete(record, "deleted_at") // ignore this, might implement soft deletes later
				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(record[key])
					delete(record, key)
				}
				r.EqualValues(tc.expectRecord, record)
			}
		})
	}
}

func TestGetAllConfigsRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []*model.AppConfigPayload
		reqData       reqData
		expectRecords []map[string]interface{}
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
				},
			},
			expectRecords: []map[string]interface{}{},
		},
		{
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "TEST", "rdev-val"),
				model.NewAppConfigPayload("testapp", "rdev", "", "TEST2", "rdev-2-val"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "TEST2", "rdev-2-stack-val"),
				model.NewAppConfigPayload("testapp", "staging", "", "TEST2", "staging-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
				},
			},
			expectRecords: []map[string]interface{}{
				{
					"app_name":    "testapp",
					"environment": "rdev",
					"key":         "TEST",
					"value":       "rdev-val",
					"source":      "environment",
					"deleted_at":  nil,
				},
				{
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
					"key":         "TEST2",
					"value":       "rdev-2-stack-val",
					"source":      "stack",
					"deleted_at":  nil,
				},
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			for _, input := range tc.seeds {
				_, err := cmd.MakeConfig(app.DB).SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, http.MethodGet, "/v1/configs", tc.reqData, r)
			count := respBody["count"].(float64)
			r.Equal(len(tc.expectRecords), int(count))

			records := respBody["records"].([]interface{})
			modifiedRecords := []map[string]interface{}{}
			for _, record := range records {
				rec := record.(map[string]interface{})

				_, createdAtPresent := rec["created_at"]
				r.Equal(true, createdAtPresent)
				_, updatedAtPresent := rec["updated_at"]
				r.Equal(true, updatedAtPresent)

				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(rec[key])
					delete(rec, key)
				}
				modifiedRecords = append(modifiedRecords, rec)
			}
			r.ElementsMatch(tc.expectRecords, modifiedRecords)
		})
	}
}

func TestCopyConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds        []*model.AppConfigPayload
		reqData      reqData
		expectRecord map[string]interface{}
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
					"key":                     "TEST2",
				},
			},
			expectRecord: nil,
		},
		{
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "TEST", "rdev-val"),
				model.NewAppConfigPayload("testapp", "rdev", "", "TEST2", "rdev-2-val"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "TEST2", "rdev-2-stack-val"),
				model.NewAppConfigPayload("testapp", "staging", "", "TEST2", "staging-val"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
					"key":                     "TEST2",
				},
			},
			expectRecord: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "staging",
				"key":         "TEST2",
				"value":       "rdev-2-stack-val",
				"deleted_at":  nil,
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			for _, input := range tc.seeds {
				_, err := cmd.MakeConfig(app.DB).SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, http.MethodPost, "/v1/config/copy", tc.reqData, r)

			if tc.expectRecord == nil {
				r.Nil(respBody["record"])
			} else {
				record := respBody["record"].(map[string]interface{})

				_, createdAtPresent := record["created_at"]
				r.Equal(true, createdAtPresent)
				_, updatedAtPresent := record["updated_at"]
				r.Equal(true, updatedAtPresent)

				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(record[key])
					delete(record, key)
				}
				r.Equal(tc.expectRecord, record)
			}
		})
	}
}

func TestCopyConfigRouteFail(t *testing.T) {
	testData := []struct {
		reqData     reqData
		failedField string
	}{
		{
			reqData: reqData{
				body: map[string]interface{}{
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
					"key":                     "TEST2",
				},
			},
			failedField: "app_name",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
					"key":                     "TEST2",
				},
			},
			failedField: "source_environment",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":           "testapp",
					"source_environment": "rdev",
					"source_stack":       "foo",
					"destination_stack":  "",
					"key":                "TEST2",
				},
			},
			failedField: "destination_environment",
		},
		{
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			failedField: "key",
		},
		{
			// copy from staging to rdev fails
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "staging",
					"source_stack":            "",
					"destination_environment": "rdev",
					"destination_stack":       "foo",
				},
			},
			failedField: "destination_environment",
		},
		{
			// copy from prod to staging fails
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "prod",
					"source_stack":            "",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			failedField: "destination_environment",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			respBody := makeInvalidRequest(app, http.MethodPost, "/v1/config/copy", tc.reqData, r)

			r.Equal(tc.failedField, respBody[0]["failed_field"])
		})
	}
}

func TestCopyDiffRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []*model.AppConfigPayload
		reqData       reqData
		expectRecords []map[string]interface{}
	}{
		{
			// no configs -> no copies
			seeds: []*model.AppConfigPayload{},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			expectRecords: []map[string]interface{}{},
		},
		{
			// config exists only for stack and no stack specified -> no copies
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			expectRecords: []map[string]interface{}{},
		},
		{
			// config exists only for env, stack specified -> env config is part of stack -> config copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			expectRecords: []map[string]interface{}{
				{
					"app_name":    "testapp",
					"environment": "staging",
					"key":         "KEY1",
					"value":       "val1",
					"deleted_at":  nil,
				},
			},
		},
		{
			// same configs in each -> no copies
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "rdev-foo-val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "staging-val1"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "foo",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			expectRecords: []map[string]interface{}{},
		},
		{
			// configs exists only in source -> configs copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "foo-val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY2", "bar-val2"),
			},
			reqData: reqData{
				body: map[string]interface{}{
					"app_name":                "testapp",
					"source_environment":      "rdev",
					"source_stack":            "",
					"destination_environment": "staging",
					"destination_stack":       "",
				},
			},
			expectRecords: []map[string]interface{}{
				{
					"app_name":    "testapp",
					"environment": "staging",
					"key":         "KEY1",
					"value":       "val1",
					"deleted_at":  nil,
				},
				{
					"app_name":    "testapp",
					"environment": "staging",
					"key":         "KEY2",
					"value":       "val2",
					"deleted_at":  nil,
				},
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(t)

			for _, input := range tc.seeds {
				_, err := cmd.MakeConfig(app.DB).SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, http.MethodPost, "/v1/config/copyDiff", tc.reqData, r)
			count := respBody["count"].(float64)
			r.Equal(len(tc.expectRecords), int(count))

			records := respBody["records"].([]interface{})
			modifiedRecords := []map[string]interface{}{}
			for _, record := range records {
				rec := record.(map[string]interface{})

				_, createdAtPresent := rec["created_at"]
				r.Equal(true, createdAtPresent)
				_, updatedAtPresent := rec["updated_at"]
				r.Equal(true, updatedAtPresent)

				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(rec[key])
					delete(rec, key)
				}
				modifiedRecords = append(modifiedRecords, rec)
			}
			r.ElementsMatch(tc.expectRecords, modifiedRecords)
		})
	}
}
