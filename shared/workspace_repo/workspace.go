package workspace_repo

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
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

type TFEMessage struct {
	Level     string `json:"@level"`
	Message   string `json:"@message"`
	Type      string `json:"type"`
	Terraform string `json:"terraform"`
	UI        string `json:"ui"`
}

type State struct {
	Version          int         `json:"version"`
	TerraformVersion string      `json:"terraform_version"`
	Serial           int64       `json:"serial"`
	Lineage          string      `json:"lineage"`
	Outputs          interface{} `json:"outputs"`
	Resources        []Resource  `json:"resources"`
}

type Resource struct {
	Module    string     `json:"module"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	IndexKey      interface{}            `json:"index_key"`
	SchemaVersion interface{}            `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
	Private       string                 `json:"private"`
	Dependencies  []string               `json:"dependencies"`
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

func (s *TFEWorkspace) GetWorkspaceName() string {
	return s.workspace.Name
}

func (s *TFEWorkspace) GetWorkspaceOrganizationName() string {
	if s.workspace.Organization == nil {
		return ""
	}
	return s.workspace.Organization.Name
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

func (s *TFEWorkspace) getCurrentRun(ctx context.Context) (*tfe.Run, error) {
	if s.currentRun != nil {
		return s.currentRun, nil
	}

	if s.GetCurrentRunID() == "" {
		return nil, errors.Errorf("fail to get current Run for %s: Run ID is empty", s.WorkspaceName())
	}

	currentRun, err := s.tfc.Runs.Read(ctx, s.GetCurrentRunID())
	if err != nil {
		return nil, errors.Wrap(err, "could not get tfe run")
	}
	s.currentRun = currentRun
	return s.currentRun, nil
}

func (s *TFEWorkspace) DiscardRun(ctx context.Context, runID string) error {
	if len(runID) == 0 {
		return errors.New("no run to discard")
	}
	return s.tfc.Runs.Discard(ctx, runID, tfe.RunDiscardOptions{
		Comment: tfe.String("cancelled by happy"),
	})
}

func (s *TFEWorkspace) GetLatestConfigVersionID(ctx context.Context) (string, error) {
	currentRun, err := s.getCurrentRun(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get the lastest ConfigVersion ID")
	}

	return currentRun.ConfigurationVersion.ID, nil
}

func (s *TFEWorkspace) Run(ctx context.Context, options ...TFERunOption) error {
	logrus.Debugf("running workspace %s ...", s.workspace.Name)
	lastConfigVersionId, err := s.GetLatestConfigVersionID(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to get latest config version id on workspace %s", s.workspace.Name)
	}
	err = s.RunConfigVersion(ctx, lastConfigVersionId, options...)
	if err != nil {
		return errors.Wrapf(err, "failed to run config version %s on workspace %s", lastConfigVersionId, s.workspace.Name)
	}

	return nil
}

func (s *TFEWorkspace) HasState(ctx context.Context) (bool, error) {
	options := tfe.StateVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 0,
			PageSize:   10,
		},
		Organization: s.GetWorkspaceOrganizationName(),
		Workspace:    s.WorkspaceName(),
	}
	list, err := s.tfc.StateVersions.List(ctx, &options)
	if err != nil {
		if errors.Is(err, tfe.ErrResourceNotFound) {
			return false, nil
		}
		return true, err
	}
	return len(list.Items) > 0, nil
}

func (s *TFEWorkspace) getVars(ctx context.Context) (map[string]map[string]*tfe.Variable, error) {
	if s.vars == nil {
		workspaceVars, err := s.tfc.Variables.List(ctx, s.GetWorkspaceId(), &tfe.VariableListOptions{})
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

func (s *TFEWorkspace) SetVars(ctx context.Context, key string, value string, description string, sensitive bool) error {
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
		_, err := s.tfc.Variables.Update(ctx, s.GetWorkspaceID(), variable.ID, options)
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
	if util.IsLocalstackMode() {
		return nil
	}
	var err error

	if s.vars == nil {
		s.vars = map[string]map[string]*tfe.Variable{}
	}
	if s.vars[category] == nil {
		s.vars[category] = map[string]*tfe.Variable{}
	}
	s.vars[category][key], err = s.tfc.Variables.Create(ctx, s.GetWorkspaceID(), options)
	return errors.Wrapf(err, "could not create TFE variable %s:%s", key, value)
}

type TFERunOption func(options *tfe.RunCreateOptions)

func DryRun(dryRun bool) TFERunOption {
	return func(options *tfe.RunCreateOptions) {
		options.ConfigurationVersion.Speculative = dryRun
		options.AutoApply = tfe.Bool(!dryRun)
	}
}

func Message(message string) TFERunOption {
	return func(options *tfe.RunCreateOptions) {
		options.Message = tfe.String(message)
	}
}

func TargetAddrs(targets []string) TFERunOption {
	return func(options *tfe.RunCreateOptions) {
		options.TargetAddrs = targets
	}
}

func (s *TFEWorkspace) RunConfigVersion(ctx context.Context, configVersionId string, opts ...TFERunOption) error {
	options := &tfe.RunCreateOptions{
		Type:      "runs",
		IsDestroy: tfe.Bool(false),
		Message:   tfe.String(fmt.Sprintf("Happy %s queued from cli", util.GetVersion().Version)),
		ConfigurationVersion: &tfe.ConfigurationVersion{
			ID:          configVersionId,
			Speculative: false,
		},
		Workspace: &tfe.Workspace{
			ID: s.GetWorkspaceID(),
		},
		TargetAddrs: []string{},
	}
	for _, opt := range opts {
		opt(options)
	}

	ws, err := s.tfc.Workspaces.ReadByID(ctx, options.Workspace.ID)
	if err != nil {
		return errors.Wrapf(err, "unable to find workspace %s, your TFE token permissions might not be sufficient", options.Workspace.ID)
	}
	logrus.Debugf("Found workspace %s", ws.Name)

	logrus.Debugf("version ID: %s, options: %+v", configVersionId, options)
	run, err := s.tfc.Runs.Create(ctx, *options)
	if err != nil {
		return errors.Wrapf(err, "could not create TFE run for workspace %s", options.Workspace.ID)
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

	diagnostics.AddTfeRunInfoUrl(ctx, s.tfc.BaseRegistryURL().Host)
	diagnostics.AddTfeRunInfoOrg(ctx, s.GetWorkspaceOrganizationName())
	diagnostics.AddTfeRunInfoWorkspace(ctx, s.GetWorkspaceName())
	diagnostics.AddTfeRunInfoRunId(ctx, s.GetCurrentRunID())
	diagnostics.PrintTfeRunLink(ctx)

	startTimestamp := time.Now()
	printedAlert := false

	var sentinelStatus tfe.RunStatus = ""
	lastStatus := sentinelStatus

	done := false

	logCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for done = false; !done; _, done = RunDoneStatuses[lastStatus] {
		if lastStatus != sentinelStatus {
			time.Sleep(5 * time.Second)
		}
		run, err := s.tfc.Runs.Read(ctx, s.GetCurrentRunID())
		if err != nil {
			return errors.Wrapf(err, "unable to get run status for run %s", s.GetCurrentRunID())
		}
		status := run.Status

		if waitOptions.Orchestrator != nil && !printedAlert && len(waitOptions.StackName) > 0 && time.Since(startTimestamp) > alertAfter {
			logrus.Warn("This apply is taking an unusually long time. Are your containers crashing?")
			// Not all services defined in config will be in the stack (e.g. if they are not deployed in this environment)
			err = waitOptions.Orchestrator.GetEvents(ctx, waitOptions.StackName, waitOptions.Services)
			if err != nil {
				logrus.Errorf("failed to get events: %s", err.Error())
			}
			printedAlert = true
		}

		if status != lastStatus {
			elapsed := time.Since(startTimestamp)
			logrus.Infof("[%s] -> [%s]: %s elapsed", lastStatus, status, units.HumanDuration(elapsed))
			lastStatus = status

			if status == tfe.RunPlanning {
				if run.Plan != nil && len(run.Plan.ID) > 0 {
					logs, err := s.tfc.Plans.Logs(logCtx, run.Plan.ID)
					if err != nil {
						logrus.Errorf("cannot retrieve logs: %s", err.Error())
					} else {
						go s.streamLogs(logCtx, logs)
					}
				}
			}

			if status == tfe.RunApplying {
				if run.Apply != nil && len(run.Apply.ID) > 0 {
					logs, err := s.tfc.Applies.Logs(logCtx, run.Apply.ID)
					if err != nil {
						logrus.Errorf("cannot retrieve logs: %s", err.Error())
					} else {
						go s.streamLogs(logCtx, logs)
					}
				}
			}

			if status == tfe.RunErrored {
				logrus.Errorf("TFE plan errored, please check the status at %s", s.GetCurrentRunUrl(ctx))
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
	// NOTE: in certain contexts
	// we don't want to show these unless specifically requested
	logfunc := logrus.Info
	if util.IsCI(ctx) {
		logfunc = logrus.Debug
	}

	logfunc("...streaming logs...")
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			logfunc("...log stream cancelled...")
			return
		default:
			bytes := scanner.Bytes()
			if len(bytes) > 0 && bytes[0] == 0x7b {
				var message TFEMessage
				err := json.Unmarshal(scanner.Bytes(), &message)
				if err == nil {
					if message.Level == "error" {
						logrus.Error(message.Message)
						continue
					}
					logfunc(message.Message)
					continue
				}
			}

			logfunc(string(bytes))
		}
	}
	if err := scanner.Err(); err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
			logrus.Errorf("...log stream error: %s...", err.Error())
			return
		}
	}
	logrus.Info("...log stream ended...")
}

// TODO: I'm not sure what this method is for
func (s *TFEWorkspace) ResetCache() {
	//s.vars = nil
	s.outputs = nil
	s.currentRun = nil
}

func (s *TFEWorkspace) GetHappyMetaRaw(ctx context.Context) ([]byte, error) {
	b := []byte{}

	vars, err := s.getVars(ctx)
	if err != nil {
		return nil, err
	}

	terraformVars, ok := vars["terraform"]
	if !ok {
		return b, nil
	}

	happyMetaVar, ok := terraformVars["happymeta_"]
	if !ok {
		return b, nil
	}

	if happyMetaVar.Sensitive {
		return nil, errors.Errorf("invalid meta var for stack %s, must not be sensitive", s.workspace.Name)
	}

	return []byte(happyMetaVar.Value), nil
}

func (s *TFEWorkspace) GetTags(ctx context.Context) (map[string]string, error) {
	tags := map[string]string{}

	vars, err := s.getVars(ctx)
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

func (s *TFEWorkspace) GetOutputs(ctx context.Context) (map[string]string, error) {
	if s.outputs != nil {
		return s.outputs, nil
	}

	s.outputs = map[string]string{}
	stateVersion, err := s.tfc.StateVersions.ReadCurrentWithOptions(ctx, s.GetWorkspaceId(), &tfe.StateVersionCurrentOptions{Include: []tfe.StateVersionIncludeOpt{"outputs"}})
	if err != nil {
		return nil, errors.Errorf("failed to get state for workspace %s", s.GetWorkspaceID())
	}

	var svOutputIDs []string
	for _, svOutput := range stateVersion.Outputs {
		svOutputIDs = append(svOutputIDs, svOutput.ID)
	}

	for _, svOutputID := range svOutputIDs {
		svOutput, err := s.tfc.StateVersionOutputs.Read(ctx, svOutputID)
		if err != nil {
			return nil, errors.Wrap(err, "could not read state version outputs")
		}

		if !svOutput.Sensitive {
			bytes, err := json.MarshalIndent(svOutput.Value, "", "\t")
			if err != nil {
				s.outputs[svOutput.Name] = fmt.Sprintf("%v", svOutput.Value)
			} else {
				s.outputs[svOutput.Name] = string(bytes)
			}
		}
	}

	return s.outputs, nil
}

func (s *TFEWorkspace) GetResources(ctx context.Context) ([]util.ManagedResource, error) {
	stateVersion, err := s.tfc.StateVersions.ReadCurrent(ctx, s.GetWorkspaceID())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state for workspace %s", s.GetWorkspaceID())
	}

	stateBytes, err := s.tfc.StateVersions.Download(ctx, stateVersion.DownloadURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download state for workspace %s, url %s", s.GetWorkspaceID(), stateVersion.DownloadURL)
	}

	var state State

	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal state for workspace %s", s.GetWorkspaceID())
	}

	resources := make([]util.ManagedResource, 0)
	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}
		if resource.Type == "validation_error" {
			continue
		}

		instances := make([]string, 0)
		id := ""

		for _, instance := range resource.Instances {
			for name, value := range instance.Attributes {
				if name == "id" {
					id = value.(string)
				}
				if name == "id" && strings.Contains(value.(string), "arn") {
					instances = append(instances, value.(string))
					break
				}
				if strings.Contains(name, "arn") {
					instances = append(instances, fmt.Sprintf("%s", value))
					break
				}
			}
		}

		if len(instances) == 0 {
			instances = append(instances, id)
		}

		ManagedResource := util.ManagedResource{
			Name:      resource.Name,
			Module:    resource.Module,
			Type:      resource.Type,
			Provider:  resource.Provider,
			Instances: instances,
			ManagedBy: "terraform",
		}

		resources = append(resources, ManagedResource)
	}
	return resources, err
}

func (s *TFEWorkspace) GetCurrentRunStatus(ctx context.Context) string {
	state, err := s.HasState(ctx)
	if err != nil || !state {
		return "no-state"
	}
	if s.currentRun == nil {
		currentRun, err := s.tfc.Runs.Read(ctx, s.workspace.CurrentRun.ID)
		if err != nil {
			logrus.Errorf("failed to get current run for workspace %s", s.WorkspaceName())
			return ""
		}
		s.currentRun = currentRun
	}
	return string(s.currentRun.Status)
}

func (s *TFEWorkspace) GetCurrentRunUrl(ctx context.Context) string {
	state, err := s.HasState(ctx)
	if err != nil || !state {
		return s.GetWorkspaceUrl()
	}

	return fmt.Sprintf("%s/runs/%s", s.GetWorkspaceUrl(), s.GetCurrentRunID())
}

// create a new ConfigurationVersion in a TFE workspace, upload the targz file to
// the new ConfigurationVersion, and finally return its ID.
func (s *TFEWorkspace) UploadVersion(ctx context.Context, targzFilePath string) (string, error) {
	logrus.WithField("workspace", s.GetWorkspaceName()).WithField("workspaceId", s.GetWorkspaceID()).WithField("org", s.GetWorkspaceOrganizationName()).Debug("Uploading configuration version")
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	autoQueueRun := false
	options := tfe.ConfigurationVersionCreateOptions{
		Type:          "configuration-versions",
		AutoQueueRuns: &autoQueueRun,
		Speculative:   &dryRun,
	}

	ws, err := s.tfc.Workspaces.ReadByID(ctx, s.GetWorkspaceID())
	if err != nil {
		return "", errors.Wrapf(err, "unable to find workspace %s, your TFE token permissions might not be sufficient", s.GetWorkspaceID())
	}
	logrus.Debugf("Found workspace %s", ws.Name)

	configVersion, err := s.tfc.ConfigurationVersions.Create(ctx, s.GetWorkspaceID(), options)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create configuration version for workspace %s (%s)", s.GetWorkspaceID(), s.GetWorkspaceName())
	}
	if err := s.tfc.ConfigurationVersions.Upload(ctx, configVersion.UploadURL, targzFilePath); err != nil {
		return "", errors.Wrapf(err, "failed to upload configuration version for workspace %s; uploadUrl=%s; targzFilePath=%s", s.GetWorkspaceID(), configVersion.UploadURL, targzFilePath)
	}

	uploaded, err := util.IntervalWithTimeout(func() (bool, error) {
		cv, err := s.tfc.ConfigurationVersions.Read(ctx, configVersion.ID)
		if err != nil {
			return false, errors.Wrapf(err, "Failed to retrieve configuration version")
		}
		if cv.Status != tfe.ConfigurationUploaded {
			return false, errors.New("configuration version not uploaded yet")
		}
		return true, nil
	}, 500*time.Millisecond, 1*time.Minute)

	if err != nil {
		return "", errors.Wrapf(err, "failed to upload configuration version")
	}

	if !*uploaded {
		return "", errors.New("failed to upload configuration version")
	}

	return configVersion.ID, nil
}

func (s *TFEWorkspace) GetWorkspaceUrl() string {
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s", s.tfc.BaseRegistryURL().Host, s.GetWorkspaceOrganizationName(), s.GetWorkspaceName())
}

func (s *TFEWorkspace) GetEndpoints(ctx context.Context) (map[string]string, error) {
	endpoints := map[string]string{}
	state, err := s.HasState(ctx)
	if err != nil {
		return endpoints, errors.Wrap(err, "unable to check if workspace had state")
	}
	if !state {
		return endpoints, nil
	}
	outputs, err := s.GetOutputs(ctx)
	if err != nil {
		return endpoints, errors.Wrap(err, "unable to get workspace outputs")
	}
	if endpoint, ok := outputs["frontend_url"]; ok {
		endpoints["FRONTEND"] = endpoint
	} else if svc_endpoints, ok := outputs["service_urls"]; ok {
		err := json.Unmarshal([]byte(svc_endpoints), &endpoints)
		if err != nil {
			return endpoints, errors.Wrap(err, "unable to decode endpoints")
		}
	} else if svc_endpoints, ok := outputs["service_endpoints"]; ok {
		err := json.Unmarshal([]byte(svc_endpoints), &endpoints)
		if err != nil {
			return endpoints, errors.Wrap(err, "unable to decode endpoints")
		}
	}
	return endpoints, nil
}
