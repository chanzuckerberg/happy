package workspace_repo

import (
	"context"
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

		if strings.Contains(r.URL.String(), "/api/v2/admin/runs") {
			fileName = fmt.Sprintf("./testdata%s.%s.json", "/api/v2/admin/runs", r.Method)
		}

		// HACK: grab the vars for generic workspace
		fileName = strings.Replace(fileName, "ws-R6X7RcX53px6vWoH", "workspace", 1)

		logrus.Warnf("filename %s", fileName)
		f, err := os.Open(fileName)
		req.NoError(err)
		_, err = io.Copy(w, f)
		req.NoError(err)

		w.WriteHeader(204)
	}))
	defer ts.Close()
	os.Setenv("TFE_TOKEN", "token")

	cf := &tfe.Config{
		Address:    ts.URL,
		Token:      "abcd1234",
		HTTPClient: ts.Client(),
	}

	client, err := tfe.NewClient(cf)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := NewWorkspaceRepo("http://example.com", "organization").WithTFEClient(client)

	_, err = repo.getToken("hostname")
	req.NoError(err)
	_, err = repo.getTfc(ctx)
	req.NoError(err)
	_, err = repo.Stacks()
	req.NoError(err)

	size, backlog, err := repo.EstimateBacklogSize(ctx)
	req.NoError(err)
	req.Equal(size, 1)
	req.NotEmpty(backlog)
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

		// HACK: grab the vars for generic workspace
		fileName = strings.Replace(fileName, "ws-R6X7RcX53px6vWoH", "workspace", 1)

		logrus.Warnf("filename %s", fileName)
		f, err := os.Open(fileName)
		req.NoError(err)
		_, err = io.Copy(w, f)
		req.NoError(err)

		w.WriteHeader(204)
	}))
	defer ts.Close()

	ctrl := gomock.NewController(t)
	mockWorkspaceRepo := NewMockWorkspaceRepoIface(ctrl)
	ctx := context.Background()

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
	ws.SetWorkspace(&tfe.Workspace{ID: "workspace", CurrentRun: &currentRun, Organization: &tfe.Organization{Name: "org"}})

	mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(&ws, nil)
	workspace, err := mockWorkspaceRepo.GetWorkspace(ctx, "workspace")
	req.NoError(err)
	_, err = workspace.GetLatestConfigVersionID(ctx)
	req.NoError(err)
	currentRunID := workspace.GetCurrentRunID()
	req.Equal("run-CZcmD7eagjhyX0vN", currentRunID)

	err = workspace.Run(ctx)
	req.NoError(err)

	err = workspace.Wait(ctx, false)
	req.NoError(err)

	_, err = workspace.UploadVersion(ctx, ".testdata/workspace/", false)
	req.NoError(err)

	status := workspace.GetCurrentRunStatus(ctx)
	req.Equal("applied", status)
	err = workspace.SetVars(ctx, "happy/app", "happy-app", "description", false)
	req.NoError(err)
	err = workspace.SetVars(ctx, "happy/app1", "happy-app", "description", false)
	req.NoError(err)

	_, err = workspace.GetOutputs(ctx)
	req.NoError(err)

	_, err = workspace.GetTags(ctx)
	req.NoError(err)

	repo := WorkspaceRepo{}
	repo.tfc = client
	repo.org = "org"
	workspace, err = repo.GetWorkspace(ctx, "workspace")
	req.NoError(err)

	hasState, err := workspace.HasState(ctx)
	req.NoError(err)
	req.True(hasState)
}
