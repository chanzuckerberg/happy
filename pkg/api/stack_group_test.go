package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy-api/pkg/model"
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
				"enabled":     true,
			},
			expectRecord: model.MakeAppStack("testapp", "rdev", "bar", true),
		},
		{
			reqBody: map[string]interface{}{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"enabled":     false,
			},
			expectRecord: model.MakeAppStack("testapp", "rdev", "bar", false),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			app := MakeTestApp(r)

			respBody := makeSuccessfulRequest(app.FiberApp, "POST", "/v1/stacks", testCase.reqBody, r)
			b, err := json.Marshal(respBody)
			r.NoError(err)
			stack := WrappedAppStack{}
			err = json.Unmarshal(b, &stack)
			r.NoError(err)

			r.Equal(testCase.expectRecord.App, stack.Record.App)
			r.Equal(testCase.expectRecord.Environment, stack.Record.Environment)
			r.Equal(testCase.expectRecord.Stack, stack.Record.Stack)
			r.Equal(testCase.expectRecord.Enabled, stack.Record.Enabled)
		})
	}
}
