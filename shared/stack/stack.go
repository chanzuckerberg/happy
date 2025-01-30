package stack

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
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Stack struct {
	Name         string
	stackService StackServiceIface
	meta         *StackMeta
	workspace    workspace_repo.Workspace
	executor     util.Executor
}

func NewStack(
	name string,
	service StackServiceIface,
) *Stack {
	return &Stack{
		Name:         name,
		stackService: service,
		executor:     util.NewDefaultExecutor(),
	}
}

func (s *Stack) WithExecutor(executor util.Executor) *Stack {
	s.executor = executor
	return s
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

func (s *Stack) Destroy(ctx context.Context, srcDir string, waitOptions options.WaitOptions, runOptions ...workspace_repo.TFERunOption) error {
	tmpDir, err := os.MkdirTemp("", "happy-destroy")
	if err != nil {
		return errors.Wrap(err, "unable to make temp directory for destroy plan")
	}
	defer os.RemoveAll(tmpDir)

	// the only file that needs to be copied over is providers.tf, versions.tf, variables.tf since the providers need
	// explicit configuration even when doing a delete. The rest of the files can be empty.
	for _, file := range []string{"providers.tf", "versions.tf", "variables.tf"} {
		_, err := os.Stat(filepath.Join(srcDir, file))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return errors.Wrap(err, "unable to stat file")
		}
		b, err := os.ReadFile(filepath.Join(srcDir, file))
		if err != nil {
			return errors.Wrapf(err, "unable to read %s", file)
		}
		err = os.WriteFile(filepath.Join(tmpDir, file), b, 0644)
		if err != nil {
			return errors.Wrapf(err, "unable to write a temporary file %s for destroy plan", file)
		}
	}

	return s.applyFromPath(ctx, tmpDir, waitOptions, runOptions...)
}

func (s *Stack) Wait(ctx context.Context, waitOptions options.WaitOptions) error {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}
	return workspace.WaitWithOptions(ctx, waitOptions)
}

func (s *Stack) applyFromPath(ctx context.Context, srcDir string, waitOptions options.WaitOptions, runOptions ...workspace_repo.TFERunOption) error {
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
		return errors.Wrap(err, "unable to get workspace")
	}

	logrus.WithField("meta_value", s.meta).Debug("Read meta from workspace")
	metaJSON, err := json.Marshal(s.meta)
	if err != nil {
		return errors.Wrap(err, "could not marshal json for stack meta")
	}
	owner := "unknown"
	if s.meta != nil {
		owner = s.meta.Owner
	}
	description := fmt.Sprintf(
		"%s - set by %s with happy CLI (%s)",
		"happy path metadata",
		owner,
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
			owner,
			util.GetVersion().Version,
		)
		err = workspace.SetVars(ctx, k, util.TagValueToString(v), description, false)
		if err != nil {
			return errors.Wrapf(err, "unable to set TFE workspace variable %s", k)
		}
	}

	if util.IsLocalstackMode() {
		module, diag := tfconfig.LoadModule(srcDir)
		if diag.HasErrors() {
			return errors.Wrap(err, "there was an issue loading the module")
		}

		cmdPath, err := s.executor.LookPath("tflocal")
		if err != nil {
			return errors.Wrap(err, "failed to locate tflocal")
		}

		_ = os.Remove(filepath.Join(srcDir, "localstack_providers_override.tf"))

		// Run 'terraform init'

		cmd := &exec.Cmd{
			Path:   cmdPath,
			Args:   []string{"tflocal", "init"},
			Dir:    srcDir,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		logrus.Infof("%s", cmd.String())
		logrus.Infof("... in %s", srcDir)
		if err := s.executor.Run(cmd); err != nil {
			return errors.Wrap(err, "failed to execute")
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
			tag := string(metaTags)
			tag = strings.ReplaceAll(tag, `\`, `\\`)
			tag = strings.ReplaceAll(tag, "'", "\\'")
			tfArgs = append(tfArgs, fmt.Sprintf("-var=happymeta_='%s'", tag))
		}

		// Run 'terraform plan' or 'terraform apply'

		cmd = &exec.Cmd{
			Path:   cmdPath,
			Args:   tfArgs,
			Dir:    srcDir,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		logrus.Infof("%s", cmd.String())
		logrus.Infof("... in %s", srcDir)
		if err := s.executor.Run(cmd); err != nil {
			return errors.Wrap(err, "failed to execute")
		}
		return nil
	}

	logrus.Debugf("will use tf bundle found at %s", srcDir)
	configVersionId, err := workspace.UploadVersion(ctx, srcDir)
	if err != nil {
		return errors.Wrap(err, "could not upload version")
	}

	// TODO should be able to use workspace.Run() here, as workspace.UploadVersion(srcDir)
	// should have generated a Run containing the Config Version Id
	runOptions = append(runOptions, workspace_repo.DryRun(dryRun))
	err = workspace.RunConfigVersion(ctx, configVersionId, runOptions...)
	if err != nil {
		return errors.Wrap(err, "could not run config version")
	}

	err = workspace.WaitWithOptions(ctx, waitOptions)
	return errors.Wrap(err, "could not wait for workspace")
}

func (s *Stack) Apply(ctx context.Context, srcDir string, waitOptions options.WaitOptions, runOptions ...workspace_repo.TFERunOption) error {
	return s.applyFromPath(ctx, srcDir, waitOptions, runOptions...)
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

func (s *Stack) GetStackInfo(ctx context.Context) (*model.AppStackResponse, error) {
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
		return nil, errors.Wrap(err, "getting the raw happy meta from TFE workspace")
	}
	meta := StackMeta{}
	err = json.Unmarshal(metaJSON, &meta)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling stack meta")
	}

	combinedTags := []string{meta.ImageTag}
	for imageTag := range meta.ImageTags {
		combinedTags = append(combinedTags, imageTag)
	}

	return &model.AppStackResponse{
		AppMetadata: *model.NewAppMetadata(meta.App, meta.Env, meta.StackName),
		StackMetadata: model.StackMetadata{
			Owner:              meta.Owner,
			Tag:                strings.Join(combinedTags, ", "),
			Outputs:            stackOutput,
			Endpoints:          endpoints,
			LastUpdated:        meta.UpdatedAt,
			GitRepo:            meta.Repo,
			GitSHA:             meta.GitSHA,
			GitBranch:          meta.GitBranch,
			TFEWorkspaceURL:    s.workspace.GetWorkspaceUrl(),
			TFEWorkspaceStatus: s.workspace.GetCurrentRunStatus(ctx),
			TFEWorkspaceRunURL: s.workspace.GetCurrentRunUrl(ctx),
		},
	}, nil
}
