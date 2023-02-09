package diagnostics

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TfeRunInfo struct {
	TfeUrl    string
	Org       string
	Workspace string
	RunId     string
}

func NewTfeRunInfo() *TfeRunInfo {
	return &TfeRunInfo{}
}

func (info *TfeRunInfo) AddTfeUrl(url string) {
	info.TfeUrl = url
}

func (info *TfeRunInfo) AddOrg(org string) {
	info.Org = org
}

func (info *TfeRunInfo) AddWorkspace(workspace string) {
	info.Workspace = workspace
}

func (info *TfeRunInfo) AddRunId(runId string) {
	info.RunId = runId
}

func (info *TfeRunInfo) CanMakeLink() bool {
	if info.TfeUrl == "" || info.Org == "" || info.Workspace == "" || info.RunId == "" {
		return false
	}
	return true
}

func (info *TfeRunInfo) MakeTfeRunLink() (string, error) {
	if info.CanMakeLink() {
		urlParts := []string{info.TfeUrl, "app", info.Org, "workspaces", info.Workspace, "runs", info.RunId}
		return strings.Join(urlParts, "/"), nil
	}
	return "", errors.New("TFE run info is incomplete and cannot form a link to the run")
}

func (info *TfeRunInfo) PrintTfeRunLink() {
	link, err := info.MakeTfeRunLink()
	if err == nil {
		logrus.Debugf("TFE run link: %s", link)
	}
}
