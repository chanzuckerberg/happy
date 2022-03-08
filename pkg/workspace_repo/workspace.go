package workspace_repo

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/docker/go-units"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const alertAfter = 300 * time.Second

// implements the Workspace interface
type TFEWorkspace struct {
	tfc          *tfe.Client
	workspace    *tfe.Workspace
	outputs      map[string]string
	vars         map[string]map[string]*tfe.Variable
	currentRun   *tfe.Run
	currentRunID string
}

// For testing purposes only
func (s *TFEWorkspace) SetClient(tfc *tfe.Client) {
	s.tfc = tfc
}

// For testing purposes only
func (s *TFEWorkspace) SetWorkspace(workspace *tfe.Workspace) {
	s.workspace = workspace
}

func (s *TFEWorkspace) GetWorkspaceID() string {
	return s.workspace.ID
}

func (s *TFEWorkspace) GetCurrentRunID() string {
	if s.currentRunID == "" {
		currentRun := s.workspace.CurrentRun
		if currentRun != nil {
			s.currentRunID = currentRun.ID
		}
	}
	return s.currentRunID
}

func (s *TFEWorkspace) getCurrentRun() (*tfe.Run, error) {
	if s.currentRun != nil {
		return s.currentRun, nil
	}

	if s.GetCurrentRunID() == "" {
		return nil, errors.Errorf("fail to get current Run for %s: Run ID is empty", s.WorkspaceName())
	}

	currentRun, err := s.tfc.Runs.Read(context.Background(), s.GetCurrentRunID())
	if err != nil {
		return nil, errors.Wrap(err, "could not get tfe run")
	}
	s.currentRun = currentRun
	return s.currentRun, nil
}

func (s *TFEWorkspace) GetLatestConfigVersionID() (string, error) {
	currentRun, err := s.getCurrentRun()
	if err != nil {
		return "", errors.Wrap(err, "failed to get the lastest ConfigVersion ID")
	}

	return currentRun.ConfigurationVersion.ID, nil
}

func (s *TFEWorkspace) Run(isDestroy bool) error {
	logrus.Infof("running workspace %s ...", s.workspace.Name)
	lastConfigVersionId, err := s.GetLatestConfigVersionID()
	if err != nil {
		return err
	}
	err = s.RunConfigVersion(lastConfigVersionId, isDestroy)
	if err != nil {
		return err
	}

	return nil
}

func (s *TFEWorkspace) getVars() (map[string]map[string]*tfe.Variable, error) {
	if s.vars == nil {
		workspaceVars, err := s.tfc.Variables.List(context.Background(), s.GetWorkspaceId(), tfe.VariableListOptions{})
		if err != nil {
			return nil, errors.Errorf("failed to get workspace vars: %v", err)
		}

		s.vars = map[string]map[string]*tfe.Variable{}
		for _, varItem := range workspaceVars.Items {
			category := string(varItem.Category)
			if _, ok := s.vars[category]; !ok {
				innerMap := map[string]*tfe.Variable{}
				s.vars[category] = innerMap
			}
			s.vars[category][varItem.Key] = varItem
		}
	}

	return s.vars, nil
}

func (s *TFEWorkspace) WorkspaceName() string {
	return s.workspace.Name
}

func (s *TFEWorkspace) SetVars(key string, value string, description string, sensitive bool) error {
	category := "terraform" // Hard-coded, not allowing setting environment vars directly
	isHCL := false

	variableMap, ok := s.vars[category]
	variable, variableOK := variableMap[key]
	update := ok && variableOK
	if update {
		logrus.WithField("existing vars", s.vars).
			WithField("tfe_workspace_id", s.GetWorkspaceID()).
			Debugf("tfe attempting to update variable %s:%s", key, value)
		options := tfe.VariableUpdateOptions{
			Type:        "vars",
			Key:         &key,
			Value:       &value,
			Description: &description,
			HCL:         &isHCL,
			Sensitive:   &sensitive,
		}
		_, err := s.tfc.Variables.Update(context.Background(), s.GetWorkspaceID(), variable.ID, options)
		return errors.Wrapf(err, "could not update TFE variable %s:%s", key, value)
	}

	logrus.WithField("existing vars", s.vars).
		WithField("tfe_workspace_id", s.GetWorkspaceID()).
		Debugf("tfe attempting to create variable %s:%s", key, value)

	options := tfe.VariableCreateOptions{
		Type:        "vars",
		Key:         &key,
		Value:       &value,
		Description: &description,
		Category:    tfe.Category(tfe.CategoryType(category)),
		HCL:         &isHCL,
		Sensitive:   &sensitive,
	}
	_, err := s.tfc.Variables.Create(context.Background(), s.GetWorkspaceID(), options)
	return errors.Wrapf(err, "could not create TFE variable %s:%s", key, value)
}

func (s *TFEWorkspace) RunConfigVersion(configVersionId string, isDestroy bool) error {
	// TODO: say who queued this or give more contextual info
	logrus.Debugf("version ID: %s, idDestroy: %t", configVersionId, isDestroy)
	msg := "Queued from happy cli"
	option := tfe.RunCreateOptions{
		Type:      "runs",
		IsDestroy: &isDestroy,
		Message:   &msg,
		ConfigurationVersion: &tfe.ConfigurationVersion{
			ID: configVersionId,
		},
		Workspace: &tfe.Workspace{
			ID: s.GetWorkspaceID(),
		},
		TargetAddrs: []string{},
	}
	run, err := s.tfc.Runs.Create(context.Background(), option)
	if err != nil {
		return err
	}
	// the run just created is the current run
	s.currentRunID = run.ID
	s.currentRun = nil
	s.outputs = nil
	return nil
}

func (s *TFEWorkspace) Wait(ctx context.Context) error {
	return s.WaitWithOptions(ctx, options.WaitOptions{})
}

func (s *TFEWorkspace) WaitWithOptions(ctx context.Context, waitOptions options.WaitOptions) error {
	RunDoneStatuses := map[tfe.RunStatus]bool{
		tfe.RunApplied:            true,
		tfe.RunDiscarded:          true,
		tfe.RunErrored:            true,
		tfe.RunCanceled:           true,
		tfe.RunPolicySoftFailed:   true,
		tfe.RunPlannedAndFinished: true,
	}

	TfeSuccessStatuses := map[tfe.RunStatus]struct{}{
		tfe.RunApplied:            {},
		tfe.RunPlannedAndFinished: {},
	}

	startTimestamp := time.Now()
	printedAlert := false

	var sentinelStatus tfe.RunStatus = ""
	lastStatus := sentinelStatus

	done := false

	ctx1, cancel := context.WithCancel(ctx)
	defer cancel()
	for done = false; !done; _, done = RunDoneStatuses[lastStatus] {
		if lastStatus != sentinelStatus {
			time.Sleep(5 * time.Second)
		}
		run, err := s.tfc.Runs.Read(context.Background(), s.GetCurrentRunID())
		if err != nil {
			return err
		}
		status := run.Status

		if waitOptions.Orchestrator != nil && !printedAlert && len(waitOptions.StackName) > 0 && time.Since(startTimestamp) > alertAfter {
			// TODO(el): A more helpful message
			logrus.Warn("This apply is taking an unusually long time. Are your containers crashing?")
			err = waitOptions.Orchestrator.GetEvents(waitOptions.StackName, waitOptions.Services)
			if err != nil {
				return err
			}
			printedAlert = true
		}

		if status != lastStatus {
			elapsed := time.Since(startTimestamp)
			logrus.Infof("[%s] -> [%s]: %s elapsed", lastStatus, status, units.HumanDuration(elapsed))
			lastStatus = status

			if status == "planning" {
				if run.Plan != nil && len(run.Plan.ID) > 0 {
					logs, err := s.tfc.Plans.Logs(context.Background(), run.Plan.ID)
					if err != nil {
						logrus.Errorf("cannot retrieve logs: %s", err.Error())
					}
					go s.streamLogs(ctx1, logs)
				}
			}

			if status == "applying" {
				if run.Apply != nil && len(run.Apply.ID) > 0 {
					logs, err := s.tfc.Applies.Logs(context.Background(), run.Apply.ID)
					if err != nil {
						logrus.Errorf("cannot retrieve logs: %s", err.Error())
					}
					go s.streamLogs(ctx1, logs)
				}
			}
		}
	}

	_, success := TfeSuccessStatuses[lastStatus]
	if !success {
		return errors.Errorf("error applying, ended in status %s", lastStatus)
	}

	return nil
}

func (s *TFEWorkspace) streamLogs(ctx context.Context, logs io.Reader) {
	logrus.Info("...streaming logs...")

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			logrus.Info("...log stream cancelled...")
			return
		default:
			logrus.Info(string(scanner.Text()))
		}
	}
	if err := scanner.Err(); err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
			logrus.Errorf("...log stream error: %s...", err.Error())
		}
	}
	logrus.Info("...log stream ended...")
}

func (s *TFEWorkspace) ResetCache() {
	s.vars = nil
	s.outputs = nil
	s.currentRun = nil
}

func (s *TFEWorkspace) GetTags() (map[string]string, error) {
	tags := map[string]string{}

	vars, err := s.getVars()
	if err != nil {
		return nil, err
	}

	terraformVars, ok := vars["terraform"]
	if !ok {
		return tags, nil
	}

	happyMetaVar, ok := terraformVars["happymeta_"]
	if !ok {
		return tags, nil
	}

	if happyMetaVar.Sensitive {
		return nil, errors.Errorf("invalid meta var for stack %s, must not be sensitive", s.workspace.Name)
	}

	// Timestamp tags come back as numeric values, and cannot be deserialized into map[string]string; code below
	// converts float64 to string, all other non-string value types will be converted.
	allTags := map[string]interface{}{}
	err = json.Unmarshal([]byte(happyMetaVar.Value), &allTags)
	for tag, value := range allTags {
		tags[tag] = util.TagValueToString(value)
	}
	return tags, errors.Wrap(err, "could not parse json")
}

func (s *TFEWorkspace) GetWorkspaceId() string {
	return s.workspace.ID
}

// For testing purposes only
func (s *TFEWorkspace) SetOutputs(outputs map[string]string) {
	s.outputs = outputs
}

func (s *TFEWorkspace) GetOutputs() (map[string]string, error) {
	if s.outputs != nil {
		return s.outputs, nil
	}

	s.outputs = map[string]string{}
	stateVersion, err := s.tfc.StateVersions.CurrentWithOptions(context.Background(), s.GetWorkspaceId(), &tfe.StateVersionCurrentOptions{Include: "outputs"})
	if err != nil {
		return nil, errors.Errorf("failed to get state for workspace %s", s.GetWorkspaceID())
	}

	var svOutputIDs []string
	for _, svOutput := range stateVersion.Outputs {
		svOutputIDs = append(svOutputIDs, svOutput.ID)
	}

	for _, svOutputID := range svOutputIDs {
		svOutput, err := s.tfc.StateVersionOutputs.Read(context.Background(), svOutputID)
		if err != nil {
			return nil, errors.Wrap(err, "could not read state version outputs")
		}

		if !svOutput.Sensitive {
			s.outputs[svOutput.Name] = svOutput.Value.(string)
		}
	}

	return s.outputs, nil
}

func (s *TFEWorkspace) GetCurrentRunStatus() string {
	if s.currentRun == nil {
		currentRun, err := s.tfc.Runs.Read(context.Background(), s.workspace.CurrentRun.ID)
		if err != nil {
			return ""
		}
		s.currentRun = currentRun
	}
	return string(s.currentRun.Status)
}

// create a new ConfigurationVersion in a TFE workspace, upload the targz file to
// the new ConfigurationVersion, and finally return its ID.
func (s *TFEWorkspace) UploadVersion(targzFilePath string) (string, error) {
	autoQueueRun := false
	options := tfe.ConfigurationVersionCreateOptions{
		Type:          "configuration-versions",
		AutoQueueRuns: &autoQueueRun,
	}
	configVersion, err := s.tfc.ConfigurationVersions.Create(context.Background(), s.GetWorkspaceID(), options)
	if err != nil {
		return "", err
	}
	if err := s.tfc.ConfigurationVersions.Upload(context.Background(), configVersion.UploadURL, targzFilePath); err != nil {
		return "", err
	}
	return configVersion.ID, nil
}
