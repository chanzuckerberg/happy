package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chanzuckerberg/happy-api/pkg/api"
	"github.com/chanzuckerberg/happy-api/pkg/cmd/config"
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func purgeTables(r *require.Assertions) {
	err := dbutil.PurgeTables()
	r.NoError(err)
}

func createRequest(method, route string, bodyMap map[string]interface{}, r *require.Assertions) *http.Request {
	body, err := json.Marshal(bodyMap)
	r.NoError(err)

	reader := bytes.NewReader(body)
	req := httptest.NewRequest(method, route, reader)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return req
}

func makeRequest(app *fiber.App, method, route string, bodyMap map[string]interface{}, r *require.Assertions) *http.Response {
	req := createRequest(method, route, bodyMap, r)
	resp, err := app.Test(req)
	r.NoError(err)
	return resp
}

func makeSuccessfulRequest(app *fiber.App, method, route string, bodyMap map[string]interface{}, r *require.Assertions) map[string]interface{} {
	resp := makeRequest(app, method, route, bodyMap, r)
	r.Equal(fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	r.NoError(err)

	jsonBody := map[string]interface{}{}
	err = json.Unmarshal(body, &jsonBody)
	r.NoError(err)

	return jsonBody
}

func makeInvalidRequest(app *fiber.App, method, route string, bodyMap map[string]interface{}, r *require.Assertions) []map[string]interface{} {
	resp := makeRequest(app, method, route, bodyMap, r)
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
		reqBody      map[string]interface{}
		expectRecord map[string]interface{}
	}{
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
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
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"key":         "TEST",
				"value":       "test-val2",
			},
			expectRecord: map[string]interface{}{
				"deleted_at":  nil,
				"app_name":    "testapp",
				"environment": "rdev",
				"key":         "TEST",
				"value":       "test-val2",
			},
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			respBody := makeSuccessfulRequest(app, "POST", "/configs", testCase.reqBody, r)

			record := respBody["record"].(map[string]interface{})
			for _, key := range []string{"id", "created_at", "updated_at"} {
				r.NotNil(record[key])
				delete(record, key)
			}

			r.EqualValues(testCase.expectRecord, record)
		})
	}
}

func TestSetConfigRouteFailsWithMissingValue(t *testing.T) {
	testData := []struct {
		reqBody     map[string]interface{}
		failedField string
	}{
		{
			reqBody: map[string]interface{}{
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
			},
			failedField: "app_name",
		},
		{
			reqBody: map[string]interface{}{
				"app_name": "testapp",
				"stack":    "bar",
				"key":      "TEST",
				"value":    "test-val",
			},
			failedField: "environment",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"value":       "test-val",
			},
			failedField: "key",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
			},
			failedField: "value",
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			respBody := makeInvalidRequest(app, "POST", "/configs", testCase.reqBody, r)

			r.Equal(testCase.failedField, respBody[0]["failed_field"])
		})
	}
}

func TestSetConfigRouteFailsWithMalformedValue(t *testing.T) {
	testData := []struct {
		reqBody      map[string]interface{}
		errorMessage string
	}{
		{
			reqBody: map[string]interface{}{
				"app_name":    13,
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.app_name of type string",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": 13,
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.environment of type string",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       13,
				"key":         "TEST",
				"value":       "test-val",
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.stack of type string",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "",
				"key":         13,
				"value":       "test-val",
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.key of type string",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "",
				"key":         "TEST",
				"value":       13,
			},
			errorMessage: "cannot unmarshal number into Go struct field AppConfigPayload.value of type string",
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			respBody := makeInvalidRequest(app, "POST", "/configs", testCase.reqBody, r)

			r.Contains(respBody[0]["message"], testCase.errorMessage)
		})
	}
}

func TestGetConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds        []*model.AppConfigPayload
		reqBody      map[string]interface{}
		expectRecord map[string]interface{}
	}{
		{
			// only env config exists, looking up by env -> returns env config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "dev", "", "TEST", "test-val"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "dev",
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
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "dev",
				"stack":       "bar",
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
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "dev",
				"stack":       "bar",
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
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.seeds {
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, "GET", "/configs/TEST", testCase.reqBody, r)
			record := respBody["record"].(map[string]interface{})
			for _, key := range []string{"id", "created_at", "updated_at"} {
				r.NotNil(record[key])
				delete(record, key)
			}

			r.EqualValues(testCase.expectRecord, record)
		})
	}
}

func TestDeleteConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []*model.AppConfigPayload
		reqBody       map[string]interface{}
		expectRecord  map[string]interface{}
		expectDeleted bool
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "foo",
			},
			expectRecord:  nil,
			expectDeleted: false,
		},
		{
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "TEST", "test-val"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "foo",
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
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.seeds {
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, "DELETE", "/configs/TEST", testCase.reqBody, r)
			deleted := respBody["deleted"].(bool)
			r.EqualValues(testCase.expectDeleted, deleted)

			if testCase.expectRecord == nil {
				r.Nil(respBody["record"])
			} else {
				record := respBody["record"].(map[string]interface{})
				for _, key := range []string{"id", "created_at", "updated_at", "deleted_at"} {
					r.NotNil(record[key])
					delete(record, key)
				}
				r.EqualValues(testCase.expectRecord, record)
			}
		})
	}
}

func TestGetAllConfigsRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []*model.AppConfigPayload
		reqBody       map[string]interface{}
		expectRecords []map[string]interface{}
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "foo",
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
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "foo",
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
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.seeds {
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, "GET", "/configs", testCase.reqBody, r)
			count := respBody["count"].(float64)
			r.Equal(len(testCase.expectRecords), int(count))

			records := respBody["records"].([]interface{})
			modifiedRecords := []map[string]interface{}{}
			for _, record := range records {
				rec := record.(map[string]interface{})
				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(rec[key])
					delete(rec, key)
				}
				modifiedRecords = append(modifiedRecords, rec)
			}
			r.ElementsMatch(testCase.expectRecords, modifiedRecords)
		})
	}
}

func TestCopyConfigRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds        []*model.AppConfigPayload
		reqBody      map[string]interface{}
		expectRecord map[string]interface{}
	}{
		{
			seeds: []*model.AppConfigPayload{},
			reqBody: map[string]interface{}{
				"app_name":                "testapp",
				"source_environment":      "rdev",
				"source_stack":            "foo",
				"destination_environment": "staging",
				"destination_stack":       "",
				"key":                     "TEST2",
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
			reqBody: map[string]interface{}{
				"app_name":                "testapp",
				"source_environment":      "rdev",
				"source_stack":            "foo",
				"destination_environment": "staging",
				"destination_stack":       "",
				"key":                     "TEST2",
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
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.seeds {
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app, "POST", "/config/copy", testCase.reqBody, r)

			if testCase.expectRecord == nil {
				r.Nil(respBody["record"])
			} else {
				record := respBody["record"].(map[string]interface{})
				for _, key := range []string{"id", "created_at", "updated_at"} {
					r.NotNil(record[key])
					delete(record, key)
				}
				r.Equal(testCase.expectRecord, record)
			}
		})
	}
}

func TestCopyConfigRouteFail(t *testing.T) {
	testData := []struct {
		reqBody     map[string]interface{}
		failedField string
	}{
		{
			reqBody: map[string]interface{}{
				"source_environment":      "rdev",
				"source_stack":            "foo",
				"destination_environment": "staging",
				"destination_stack":       "",
				"key":                     "TEST2",
			},
			failedField: "app_name",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":                "testapp",
				"source_stack":            "foo",
				"destination_environment": "staging",
				"destination_stack":       "",
				"key":                     "TEST2",
			},
			failedField: "source_environment",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":           "testapp",
				"source_environment": "rdev",
				"source_stack":       "foo",
				"destination_stack":  "",
				"key":                "TEST2",
			},
			failedField: "destination_environment",
		},
		{
			reqBody: map[string]interface{}{
				"app_name":                "testapp",
				"source_environment":      "rdev",
				"source_stack":            "foo",
				"destination_environment": "staging",
				"destination_stack":       "",
			},
			failedField: "key",
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			app, err := api.MakeApp()
			r.NoError(err)
			defer purgeTables(r)

			respBody := makeInvalidRequest(app, "POST", "/config/copy", testCase.reqBody, r)

			r.Equal(testCase.failedField, respBody[0]["failed_field"])
		})
	}
}
