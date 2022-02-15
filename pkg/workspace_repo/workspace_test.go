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
	req := require.New(t)
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
			_, err := w.Write([]byte(response))
			req.NoError(err)
		}
		if r.URL.String() == "/api/v2/runs/run-CZcmD7eagjhyX0vN" {
			response = strings.Replace(response, "$STATUS", "applied", 1)
			_, err := w.Write([]byte(response))
			req.NoError(err)
		}
		if r.URL.String() == "/api/v2/workspaces/workspace/configuration-versions" {
			response = `{
				"data": {
				  "id": "cv-ntv3HbhJqvFzamy7",
				  "type": "configuration-versions",
				  "attributes": {
					"error": null,
					"error-message": null,
					"source": "gitlab",
					"speculative":false,
					"status": "uploaded",
					"status-timestamps": {}
				  },
				  "relationships": {
					"ingress-attributes": {
					  "data": {
						"id": "ia-i4MrTxmQXYxH2nYD",
						"type": "ingress-attributes"
					  },
					  "links": {
						"related":
						  "/api/v2/configuration-versions/cv-ntv3HbhJqvFzamy7/ingress-attributes"
					  }
					}
				  },
				  "links": {
					"self": "/api/v2/configuration-versions/cv-ntv3HbhJqvFzamy7"
				  }
				}
			  }`

			_, err := w.Write([]byte(response))
			req.NoError(err)
		}
		w.WriteHeader(204)
	}))
	defer ts.Close()

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
	currentRun := tfe.Run{ID: "run-CZcmD7eagjhyX0vN", ConfigurationVersion: &tfe.ConfigurationVersion{ID: "123"}}
	ws.SetClient(client)
	ws.SetWorkspace(&tfe.Workspace{ID: "workspace", CurrentRun: &currentRun})
	ws.SetCurrentRun(&currentRun)
	mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(&ws, nil)
	workspace, err := mockWorkspaceRepo.GetWorkspace("workspace")
	req.NoError(err)
	_, err = workspace.GetLatestConfigVersionID()
	req.NoError(err)
	currentRunID := workspace.GetCurrentRunID()
	req.Equal("run-CZcmD7eagjhyX0vN", currentRunID)

	err = workspace.Run(false)
	req.NoError(err)

	err = workspace.Wait()
	req.NoError(err)

	_, err = workspace.UploadVersion("../config/testdata/")
	req.NoError(err)
}
