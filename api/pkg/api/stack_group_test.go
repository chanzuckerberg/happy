package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy-api/pkg/cmd"
	"github.com/chanzuckerberg/happy-shared/model"
	"github.com/stretchr/testify/require"
)

func TestCreateStackRouteSucceed(t *testing.T) {
	testData := []struct {
		reqBody      map[string]interface{}
		expectRecord model.AppStack
	}{
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
			},
			expectRecord: model.MakeAppStack("testapp", "rdev", "bar"),
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
			},
			expectRecord: model.MakeAppStack("testapp", "rdev", "bar"),
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(r)

			respBody := makeSuccessfulRequest(app.FiberApp, "POST", "/v1/stacklistItems", tc.reqBody, r)
			b, err := json.Marshal(respBody)
			r.NoError(err)
			stack := WrappedAppStack{}
			err = json.Unmarshal(b, &stack)
			r.NoError(err)

			r.Equal(tc.expectRecord.App, stack.Record.App)
			r.Equal(tc.expectRecord.Environment, stack.Record.Environment)
			r.Equal(tc.expectRecord.Stack, stack.Record.Stack)
		})
	}
}

func TestGetStacklistRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds         []model.AppStackPayload
		reqBody       map[string]interface{}
		expectRecords []map[string]interface{}
	}{
		{
			// nothing exists -> no records returned
			seeds: []model.AppStackPayload{},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
			},
			expectRecords: []map[string]interface{}{},
		},
		{
			// when record exists -> it is returned
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "bar"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
			},
			expectRecords: []map[string]interface{}{
				{
					"deleted_at":  nil,
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "bar",
				},
			},
		},
		{
			// when many records exist -> matching records are returned
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "foo"),
				model.MakeAppStackPayload("testapp", "rdev", "bar"),
				model.MakeAppStackPayload("testapp", "staging", "foo"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
			},
			expectRecords: []map[string]interface{}{
				{
					"deleted_at":  nil,
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "foo",
				},
				{
					"deleted_at":  nil,
					"app_name":    "testapp",
					"environment": "rdev",
					"stack":       "bar",
				},
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(r)

			for _, input := range tc.seeds {
				_, err := cmd.MakeStack(app.DB).CreateOrUpdateAppStack(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app.FiberApp, "GET", "/v1/stacklistItems/", tc.reqBody, r)

			count := respBody["count"].(float64)
			r.Equal(len(tc.expectRecords), int(count))

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
			r.ElementsMatch(tc.expectRecords, modifiedRecords)
		})
	}
}

func TestDeleteStacklistItemRouteSucceed(t *testing.T) {
	testData := []struct {
		seeds        []model.AppStackPayload
		reqBody      map[string]interface{}
		expectRecord map[string]interface{}
	}{
		{
			// when record exists -> it is deleted
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "bar"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
			},
			expectRecord: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
			},
		},
		{
			// when record does not exist -> nothing is deleted
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "foo"),
			},
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
			},
			expectRecord: nil,
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(r)

			for _, input := range tc.seeds {
				_, err := cmd.MakeStack(app.DB).CreateOrUpdateAppStack(input)
				r.NoError(err)
			}

			respBody := makeSuccessfulRequest(app.FiberApp, "DELETE", "/v1/stacklistItems/", tc.reqBody, r)

			if tc.expectRecord == nil {
				r.Nil(respBody["record"])
			} else {
				record := respBody["record"].(map[string]interface{})
				for _, key := range []string{"id", "created_at", "updated_at", "deleted_at"} {
					r.NotNil(record[key])
					delete(record, key)
				}
				fmt.Println("testCase.expectRecord: ", tc.expectRecord)
				fmt.Println("record: ", record)
				r.EqualValues(tc.expectRecord, record)
			}
		})
	}
}
