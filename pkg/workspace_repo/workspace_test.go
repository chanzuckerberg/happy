package workspace_repo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/sirupsen/logrus"
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
		logrus.Infof("%s %s\n", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-RateLimit-Limit", "30")
		w.Header().Set("TFP-API-Version", "34.21.9")
		if r.URL.String() == "/api/v2/ping" {
			w.WriteHeader(204)
			return
		}

		fileName := fmt.Sprintf("./testdata%s.%s.json", r.URL.String(), r.Method)
		if strings.Contains(r.URL.String(), "/api/v2/state-version-outputs/") {
			fileName = fmt.Sprintf("./testdata%s.%s.json", "/api/v2/state-version-outputs", r.Method)
		}
		f, err := os.Open(fileName)
		req.NoError(err)
		_, err = io.Copy(w, f)
		req.NoError(err)

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

	status := workspace.GetCurrentRunStatus()
	req.Equal("applied", status)
	err = workspace.SetVars("happy/app", "happy-app", "description", false)
	req.NoError(err)
	err = workspace.SetVars("happy/app1", "happy-app", "description", false)
	req.NoError(err)

	_, err = workspace.GetOutputs()
	req.NoError(err)

	_, err = workspace.GetTags()
	req.NoError(err)

	repo := WorkspaceRepo{}
	repo.tfc = client
	repo.org = "org"
	_, err = repo.GetWorkspace("workspace")
	req.NoError(err)
}
