package workspace_repo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceRepo(t *testing.T) {
	r := require.New(t)
	os.Setenv("TFE_TOKEN", "token")
	repo, err := NewWorkspaceRepo("https://repo.com", "organization")
	r.NoError(err)
	_, err = repo.getToken("hostname")
	r.NoError(err)
	_, err = repo.getTfc()
	r.NoError(err)
	_, err = repo.Stacks()
	r.NoError(err)
}

func TestWorkspace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-RateLimit-Limit", "30")
		w.Header().Set("TFP-API-Version", "34.21.9")
		if r.URL.String() == "/api/v2/ping" {
			w.WriteHeader(204)
			return
		}

		response := `{
			"data": {
			  "id": "run-CZcmD7eagjhyX0vN",
			  "type": "runs",
			  "attributes": {
				"actions": {
				  "is-cancelable": true,
				  "is-confirmable": false,
				  "is-discardable": false,
				  "is-force-cancelable": false
				},
				"canceled-at": null,
				"created-at": "2021-05-24T07:38:04.171Z",
				"has-changes": false,
				"auto-apply": false,
				"is-destroy": false,
				"message": "Custom message",
				"plan-only": false,
				"source": "tfe-api",
				"status-timestamps": {
				  "plan-queueable-at": "2021-05-24T07:38:04+00:00"
				},
				"status": "$STATUS",
				"trigger-reason": "manual",
				"target-addrs": null,
				"permissions": {
				  "can-apply": true,
				  "can-cancel": true,
				  "can-comment": true,
				  "can-discard": true,
				  "can-force-execute": true,
				  "can-force-cancel": true,
				  "can-override-policy-check": true
				},
				"refresh": false,
				"refresh-only": false,
				"replace-addrs": null,
				"variables": []
			  },
			  "relationships": {
				
				
			  },
			  "links": {
				"self": "/api/v2/runs/run-CZcmD7eagjhyX0vN"
			  }
			}
		  }`

		if r.URL.String() == "/api/v2/runs" {
			response = strings.Replace(response, "$STATUS", "pending", 1)
			w.Write([]byte(response))
		}
		if strings.Contains(r.URL.String(), "/api/v2/runs/run-") {
			response = strings.Replace(response, "$STATUS", "applied", 1)
			w.Write([]byte(response))
		}
		w.WriteHeader(204)
	}))
	defer ts.Close()

	r := require.New(t)
	ctrl := gomock.NewController(t)
	mockWorkspaceRepo := NewMockWorkspaceRepoIface(ctrl)

	config := &tfe.Config{
		Address:    ts.URL,
		Token:      "abcd1234",
		HTTPClient: ts.Client(),
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	ws := TFEWorkspace{}
	ws.SetClient(client)
	ws.SetWorkspace(&tfe.Workspace{})
	ws.SetCurrentRun(&tfe.Run{ConfigurationVersion: &tfe.ConfigurationVersion{ID: "123"}})
	mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(&ws, nil)
	workspace, err := mockWorkspaceRepo.GetWorkspace("workspace")
	r.NoError(err)
	_, err = workspace.GetLatestConfigVersionID()
	r.NoError(err)
	currentRunID := workspace.GetCurrentRunID()
	r.Equal("", currentRunID)

	err = workspace.Run(false)
	r.NoError(err)

	err = workspace.Wait()
	r.NoError(err)
}
