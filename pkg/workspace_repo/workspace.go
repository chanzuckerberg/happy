package workspace_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chanzuckerberg/happy/pkg/options"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

const alertAfter time.Duration = 300 * time.Second

// implements the Workspace interface
type TFEWorkspace struct {
	tfc          *tfe.Client
	workspace    *tfe.Workspace
	outputs      map[string]string
	vars         map[string]map[string]*tfe.Variable
	currentRun   *tfe.Run
	currentRunID string
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
	if s.currentRun == nil {
		if s.GetCurrentRunID() != "" {
			currentRun, err := s.tfc.Runs.Read(context.Background(), s.GetCurrentRunID())
			if err != nil {
				return nil, err
			}
			s.currentRun = currentRun
		} else {
			return nil, errors.Errorf("fail to get current Run for %s: Run ID is empty", s.WorkspaceName())
		}
	}
	return s.currentRun, nil
}

func (s *TFEWorkspace) GetLatestConfigVersionID() (string, error) {
	currentRun, err := s.getCurrentRun()
	if err != nil {
		return "", errors.Errorf("fail to get the lastest ConfigVersion ID: %s", err)
	}

	return currentRun.ConfigurationVersion.ID, nil
}

func (s *TFEWorkspace) Run(isDestroy bool) error {
	fmt.Printf("Runing workspace %s ...\n", s.workspace.Name)
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

	if variableMap, ok := s.vars[category]; ok {
		if variable, ok := variableMap[key]; ok {
			options := tfe.VariableUpdateOptions{
				Type:        "vars",
				Key:         &key,
				Value:       &value,
				Description: &description,
				HCL:         &isHCL,
				Sensitive:   &sensitive,
			}
			_, err := s.tfc.Variables.Update(context.Background(), s.GetWorkspaceID(), variable.ID, options)
			if err != nil {
				return err
			}
		}
	} else {
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
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TFEWorkspace) RunConfigVersion(configVersionId string, isDestroy bool) error {
	// TODO: say who queued this or give more contextual info
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

func (s *TFEWorkspace) Wait() error {
	return s.WaitWithOptions(options.WaitOptions{})
}

func (s *TFEWorkspace) WaitWithOptions(waitOptions options.WaitOptions) error {
	RUN_DONE_STATUSES := map[tfe.RunStatus]bool{
		tfe.RunApplied:          true,
		tfe.RunDiscarded:        true,
		tfe.RunErrored:          true,
		tfe.RunCanceled:         true,
		tfe.RunPolicySoftFailed: true,
	}

	startTimestamp := time.Now()
	printedAlert := false

	var sentinelStatus tfe.RunStatus = ""
	lastStatus := sentinelStatus

	for done := false; !done; _, done = RUN_DONE_STATUSES[lastStatus] {
		if lastStatus != sentinelStatus {
			time.Sleep(5 * time.Second)
		}
		run, err := s.tfc.Runs.Read(context.Background(), s.GetCurrentRunID())
		if err != nil {
			return err
		}
		status := run.Status

		if waitOptions.Orchestrator != nil && !printedAlert && len(waitOptions.StackName) > 0 && time.Since(startTimestamp) > alertAfter {
			log.Println("This apply is taking an unusually long time. Are your containers crashing?")
			err = waitOptions.Orchestrator.GetEvents(waitOptions.StackName, waitOptions.Services)
			if err != nil {
				return err
			}
			printedAlert = true
		}

		if status != lastStatus {
			fmt.Printf("[timestamp] - %s\n", status)
			lastStatus = status
			startTimestamp = time.Now()
		}
	}

	if lastStatus != tfe.RunApplied {
		return errors.Errorf("error applying, ended in status %s", lastStatus)
	}

	return nil
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
		return nil, errors.New("invalid meta var for stack {self.stack_name}, must not be sensitive")
	}

	err = json.Unmarshal([]byte(happyMetaVar.Value), &tags)
	return tags, errors.Wrap(err, "could not parse json")
}

func (s *TFEWorkspace) GetWorkspaceId() string {
	return s.workspace.ID
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
			return nil, err
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
