package stack_mgr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StackInfo struct {
	Name        string            `json:",omitempty"`
	Owner       string            `json:",omitempty"`
	Tag         string            `json:",omitempty"`
	Status      string            `json:",omitempty"`
	LastUpdated string            `json:",omitempty"`
	Message     string            `json:",omitempty"`
	Outputs     map[string]string `json:",omitempty"`
	Endpoints   map[string]string `json:",omitempty"`
	Repo        string            `json:",omitempty"`
	App         string            `json:",omitempty"`
	GitSHA      string            `json:",omitempty"`
	GitBranch   string            `json:",omitempty"`
}

type Stack struct {
	Name string

	stackService StackServiceIface

	meta      *StackMeta
	workspace workspace_repo.Workspace
}

func NewStack(
	name string,
	service StackServiceIface,
) *Stack {
	return &Stack{
		Name:         name,
		stackService: service,
	}
}

func (s *Stack) getWorkspace(ctx context.Context) (workspace_repo.Workspace, error) {
	var err error
	s.workspace, err = s.stackService.GetStackWorkspace(ctx, s.Name)
	return s.workspace, errors.Wrapf(err, "failed to get workspace for stack %s", s.Name)
}

func (s *Stack) GetOutputs(ctx context.Context) (map[string]string, error) {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	outputs, err := workspace.GetOutputs(ctx)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (s *Stack) GetResources(ctx context.Context) ([]util.ManagedResource, error) {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := workspace.GetResources(ctx)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (s *Stack) GetStatus(ctx context.Context) string {
	if s.workspace != nil {
		return s.workspace.GetCurrentRunStatus(ctx)
	}
	return ""
}

func (s *Stack) WithMeta(meta *StackMeta) *Stack {
	logrus.WithField("meta", meta).Debug("setting meta")
	s.meta = meta
	return s
}

func (s *Stack) Destroy(ctx context.Context) error {
	return s.PlanDestroy(ctx)
}

func (s *Stack) PlanDestroy(ctx context.Context, opts ...workspace_repo.TFERunOption) error {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Destroy")
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}

	versionId, err := workspace.GetLatestConfigVersionID(ctx)
	if err != nil {
		return err
	}

	opts = append(opts, workspace_repo.DestroyPlan(), workspace_repo.DryRun(dryRun))

	err = workspace.RunConfigVersion(ctx, versionId,
		opts...,
	)
	if err != nil {
		return err
	}
	currentRunID := workspace.GetCurrentRunID()

	err = workspace.Wait(ctx)
	if err != nil {
		return err
	}

	if dryRun {
		err = workspace.DiscardRun(ctx, currentRunID)
	}
	return err
}

func (s *Stack) Wait(ctx context.Context, waitOptions options.WaitOptions) error {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}
	return workspace.WaitWithOptions(ctx, waitOptions)
}

func (s *Stack) Apply(ctx context.Context, waitOptions options.WaitOptions, runOptions ...workspace_repo.TFERunOption) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Apply")
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		logrus.Debug()
		logrus.Debugf("planning stack %s...", s.Name)
	} else {
		logrus.Debug()
		logrus.Debugf("applying stack %s...", s.Name)
	}

	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}

	logrus.WithField("meta_value", s.meta).Debug("Read meta from workspace")
	metaJSON, err := json.Marshal(s.meta)
	if err != nil {
		return errors.Wrap(err, "could not marshal json for stack meta")
	}
	description := fmt.Sprintf(
		"%s - set by %s with happy CLI (%s)",
		"happy path metadata",
		s.meta.Owner,
		util.GetVersion().Version,
	)
	err = workspace.SetVars(ctx, "happymeta_", string(metaJSON), description, false)
	if err != nil {
		return errors.Wrap(err, "unable to set TFE workspace variable happymeta_")
	}

	metaKeys := map[string]any{}
	err = json.Unmarshal(metaJSON, &metaKeys)
	if err != nil {
		return errors.Wrap(err, "unable to json unmarshal meta keys")
	}
	for k, v := range metaKeys {
		description = fmt.Sprintf(
			"%s - set by %s with happy CLI (%s)",
			k,
			s.meta.Owner,
			util.GetVersion().Version,
		)
		err = workspace.SetVars(ctx, k, util.TagValueToString(v), description, false)
		if err != nil {
			return errors.Wrapf(err, "unable to set TFE workspace variable %s", k)
		}
	}

	tfDirPath := s.stackService.GetConfig().TerraformDirectory()

	happyProjectRoot := s.stackService.GetConfig().GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)

	if util.IsLocalstackMode() {
		module, diag := tfconfig.LoadModule(srcDir)
		if diag.HasErrors() {
			return errors.Wrap(err, "There was an issue loading the module")
		}

		// Clear out any prior state... For now. Every stack has to have its own

		// _ = os.Remove(filepath.Join(srcDir, "terraform.tfstate"))
		// _ = os.Remove(filepath.Join(srcDir, "terraform.tfstate.backup"))
		_ = os.Remove(filepath.Join(srcDir, "localstack_providers_override.tf"))

		// Run 'terraform init'

		cmd := exec.CommandContext(ctx, "tflocal", "init")
		cmd.Dir = srcDir
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		logrus.Infof("%s", cmd.String())
		logrus.Infof("... in %s", srcDir)
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "error executing tflocal init")
		}

		command := "apply"
		if dryRun {
			command = "plan"
		}
		tfArgs := []string{"tflocal", command}
		if !dryRun {
			tfArgs = append(tfArgs, "-auto-approve")
		}

		// Every stack has to have its own state file.
		tfArgs = append(tfArgs, fmt.Sprintf("-state=%s.tfstate", s.Name))
		tfArgs = append(tfArgs, "-lock=false")

		for param, value := range metaKeys {
			if _, ok := module.Variables[param]; ok {
				tfArgs = append(tfArgs, fmt.Sprintf("-var=%s=%s", param, value))
			}
		}
		metaTags, err := json.Marshal(s.meta)
		if err != nil {
			return errors.Wrap(err, "could not marshal json")
		}
		if _, ok := module.Variables["happymeta_"]; ok {
			tfArgs = append(tfArgs, fmt.Sprintf("-var=happymeta_='%s'", string(metaTags)))
		}

		// Run 'terraform plan' or 'terraform apply'

		cmd = exec.CommandContext(ctx, "tflocal", tfArgs...)
		cmd.Dir = srcDir
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Infof("%s", cmd.String())
		logrus.Infof("... in %s", srcDir)
		return errors.Wrap(cmd.Run(), "failed to execute")
	}

	logrus.Debugf("will use tf bundle found at %s", srcDir)

	tempFile, err := os.CreateTemp("", "happy_tfe.*.tar.gz")
	if err != nil {
		return errors.Wrap(err, "could not create temporary file")
	}
	defer os.Remove(tempFile.Name())
	logrus.Debugf("tarzipping file (%s) ...", tempFile.Name())
	err = util.TarDir(srcDir, tempFile)
	if err != nil {
		return err
	}

	configVersionId, err := workspace.UploadVersion(ctx, srcDir)
	if err != nil {
		return errors.Wrap(err, "could not upload version")
	}

	// TODO should be able to use workspace.Run() here, as workspace.UploadVersion(srcDir)
	// should have generated a Run containing the Config Version Id
	runOptions = append(runOptions, workspace_repo.DryRun(dryRun))
	err = workspace.RunConfigVersion(ctx, configVersionId, runOptions...)
	if err != nil {
		return err
	}

	return workspace.WaitWithOptions(ctx, waitOptions)
}

func (s *Stack) PrintOutputs(ctx context.Context) {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return
	}

	logrus.Info("Module Outputs --")
	stackOutput, err := s.GetOutputs(ctx)
	if err != nil {
		logrus.Errorf("Failed to get output for stack %s: %s", s.Name, err.Error())
		return
	}
	for k, v := range stackOutput {
		logrus.Printf("%s: %s", k, v)
	}
}

func (s *Stack) GetStackInfo(ctx context.Context, name string) (*StackInfo, error) {
	stackOutput, err := s.GetOutputs(ctx)
	if err != nil {
		return nil, err
	}
	endpoints, err := s.workspace.GetEndpoints(ctx)
	if err != nil {
		return nil, err
	}

	metaJSON, err := s.workspace.GetHappyMetaRaw(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get the raw happy meta from TFE workspace")
	}
	meta := StackMeta{}
	err = json.Unmarshal(metaJSON, &meta)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal stack meta")
	}
	// TODO: only here until people upgrade their CLIS. remove this in a few weeks
	metaLegacy := StackMetaLegacy{}
	err = json.Unmarshal(metaJSON, &metaLegacy)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal legacy stack meta")
	}
	err = meta.Merge(metaLegacy)
	if err != nil {
		return nil, errors.Wrap(err, "could not merge stack meta")
	}

	combinedTags := []string{meta.ImageTag}
	for imageTag := range meta.ImageTags {
		combinedTags = append(combinedTags, imageTag)
	}

	return &StackInfo{
		Name:        meta.StackName,
		Owner:       meta.Owner,
		Tag:         strings.Join(combinedTags, ", "),
		Status:      s.GetStatus(ctx),
		Outputs:     stackOutput,
		Endpoints:   endpoints,
		LastUpdated: meta.UpdatedAt,
		Repo:        meta.Repo,
		App:         meta.App,
		GitSHA:      meta.GitSHA,
		GitBranch:   meta.GitBranch,
	}, nil
}
