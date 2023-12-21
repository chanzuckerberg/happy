package api

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// func sendPostRequest(app *APIApplication, route string, r *require.Assertions) *http.Response {
// 	svr := httptest.NewServer(app.mux)
// 	defer svr.Close()
// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", svr.URL, route), nil)
// 	r.NoError(err)

// 	client := http.DefaultClient
// 	resp, err := client.Do(req)
// 	r.NoError(err)

// 	return resp
// }

func TestSetConfigV2RouteSucceed(t *testing.T) {
	testData := []struct {
		reqData      reqData
		expectRecord map[string]string
	}{
		{
			reqData: reqData{
				queryParams: map[string]string{
					"app_name":       "testapp",
					"environment":    "rdev",
					"stack":          "bar",
					"key":            "TEST",
					"value":          "test-val",
					"aws_profile":    "test-profile",
					"aws_region":     "us-west-2",
					"k8s_namespace":  "test-namespace",
					"k8s_cluster_id": "test-cluster-id",
				},
				headers: map[string]string{
					"X-Aws-Access-Key-Id":     b64.StdEncoding.EncodeToString([]byte("test-access-key-id")),
					"X-Aws-Secret-Access-Key": b64.StdEncoding.EncodeToString([]byte("test-secret-access-key")),
					"X-Aws-Session-Token":     b64.StdEncoding.EncodeToString([]byte("test-session-token")),
				},
			},
			expectRecord: map[string]string{
				"app_name":    "testapp",
				"environment": "rdev",
				"stack":       "bar",
				"key":         "TEST",
				"value":       "test-val",
			},
		},
		{
			reqData: reqData{
				queryParams: map[string]string{
					"app_name":       "testapp",
					"environment":    "rdev",
					"key":            "TEST",
					"value":          "test-val2",
					"aws_profile":    "test-profile",
					"aws_region":     "us-west-2",
					"k8s_namespace":  "test-namespace",
					"k8s_cluster_id": "test-cluster-id",
				},
				headers: map[string]string{
					"X-Aws-Access-Key-Id":     b64.StdEncoding.EncodeToString([]byte("test-access-key-id")),
					"X-Aws-Secret-Access-Key": b64.StdEncoding.EncodeToString([]byte("test-secret-access-key")),
					"X-Aws-Session-Token":     b64.StdEncoding.EncodeToString([]byte("test-session-token")),
				},
			},
			expectRecord: map[string]string{
				"app_name":    "testapp",
				"environment": "rdev",
				"key":         "TEST",
				"value":       "test-val2",
			},
		},
		{
			// test that special characters are standardized
			reqData: reqData{
				queryParams: map[string]string{
					"app_name":       "testapp",
					"environment":    "rdev",
					"key":            "TEST-2*()$",
					"value":          "test-val2",
					"aws_profile":    "test-profile",
					"aws_region":     "us-west-2",
					"k8s_namespace":  "test-namespace",
					"k8s_cluster_id": "test-cluster-id",
				},
				headers: map[string]string{
					"X-Aws-Access-Key-Id":     b64.StdEncoding.EncodeToString([]byte("test-access-key-id")),
					"X-Aws-Secret-Access-Key": b64.StdEncoding.EncodeToString([]byte("test-secret-access-key")),
					"X-Aws-Session-Token":     b64.StdEncoding.EncodeToString([]byte("test-session-token")),
				},
			},
			expectRecord: map[string]string{
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

			// respBody := makeSuccessfulRequest(app, http.MethodPost, "/v2/configs", tc.reqBody, r)
			// values := url.Values{}
			// for k, v := range tc.reqBody {
			// 	values.Add(k, v)
			// }
			// query := values.Encode()
			resp := makeSuccessfulRequest(app, http.MethodPost, "/v2/app-configs", tc.reqData, r)

			// record := respBody["record"].(map[string]interface{})

			// _, createdAtPresent := record["created_at"]
			// r.Equal(true, createdAtPresent)
			// _, updatedAtPresent := record["updated_at"]
			// r.Equal(true, updatedAtPresent)

			// for _, key := range []string{"id", "created_at", "updated_at"} {
			// 	r.NotNil(record[key])
			// 	delete(record, key)
			// }

			r.EqualValues(tc.expectRecord, resp)
		})
	}
}
