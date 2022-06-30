package stack_mgr

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StackIface interface {
	GetName() string
	SetMeta(meta *StackMeta)
	Meta() (*StackMeta, error)
	GetStatus() string
	GetOutputs(ctx context.Context) (map[string]string, error)
	Wait(ctx context.Context, waitOptions options.WaitOptions)
	Plan(ctx context.Context, waitOptions options.WaitOptions, dryRun util.DryRunType) error
	PlanDestroy(ctx context.Context, dryRun util.DryRunType) error
	Destroy(ctx context.Context) error
	PrintOutputs()
}

type Stack struct {
	stackName string

	stackService StackServiceIface
	dirProcessor util.DirProcessor

	meta      *StackMeta
	workspace workspace_repo.Workspace
}

func NewStack(
	name string,
	service StackServiceIface,
	dirProcessor util.DirProcessor,
) *Stack {
	return &Stack{
		stackName:    name,
		stackService: service,
		dirProcessor: dirProcessor,
	}
}

func (s *Stack) GetName() string {
	return s.stackName
}

func (s *Stack) getWorkspace(ctx context.Context) (workspace_repo.Workspace, error) {
	if s.workspace == nil {
		workspace, err := s.stackService.GetStackWorkspace(ctx, s.stackName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get workspace for stack %s", s.stackName)
		}
		s.workspace = workspace
	}

	return s.workspace, nil
}

func (s *Stack) GetOutputs(ctx context.Context) (map[string]string, error) {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	outputs, err := workspace.GetOutputs()
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (s *Stack) GetStatus() string {
	if s.workspace != nil {
		return s.workspace.GetCurrentRunStatus()
	}
	return ""
}

func (s *Stack) WithMeta(meta *StackMeta) *Stack {
	logrus.WithField("meta", meta).Debug("setting meta")
	s.meta = meta
	return s
}

func (s *Stack) Meta(ctx context.Context) (*StackMeta, error) {
	if s.meta == nil {
		s.meta = s.stackService.NewStackMeta(s.stackName)

		// update tags of meta with those from the backing workspace
		workspace, err := s.getWorkspace(ctx)
		if err != nil {
			return nil, err
		}

		tags, err := workspace.GetTags()
		if err != nil {
			return nil, err
		}

		// FIXME TODO(el): why is this? don't we want to set these?
		if len(tags) == 0 {
			tags = map[string]string{
				"happy/meta/owner":    unknown,
				"happy/meta/imagetag": unknown,
			}
		}

		err = s.meta.Load(tags)
		if err != nil {
			return nil, err
		}
	}
	return s.meta, nil
}

func (s *Stack) Destroy(ctx context.Context) error {
	return s.PlanDestroy(ctx, false)
}

func (s *Stack) PlanDestroy(ctx context.Context, dryRun util.DryRunType) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Destroy")
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}

	versionId, err := workspace.GetLatestConfigVersionID()

	if err != nil {
		return err
	}

	isDestroy := true
	err = workspace.RunConfigVersion(versionId, isDestroy, dryRun)
	if err != nil {
		return err
	}
	currentRunID := workspace.GetCurrentRunID()

	err = workspace.Wait(ctx, dryRun)
	if err != nil {
		return err
	}

	if dryRun {
		err = workspace.DiscardRun(ctx, currentRunID)
	}
	return err
}

func (s *Stack) Wait(ctx context.Context, waitOptions options.WaitOptions, dryRun util.DryRunType) error {
	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}
	return workspace.WaitWithOptions(ctx, waitOptions, dryRun)
}

func (s *Stack) Apply(ctx context.Context, waitOptions options.WaitOptions) error {
	return s.Plan(ctx, waitOptions, false)
}

func (s *Stack) Plan(ctx context.Context, waitOptions options.WaitOptions, dryRun util.DryRunType) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Apply")
	if dryRun {
		logrus.Info()
		logrus.Infof("planning stack %s...", s.stackName)
	} else {
		logrus.Info()
		logrus.Infof("applying stack %s...", s.stackName)
	}

	workspace, err := s.getWorkspace(ctx)
	if err != nil {
		return err
	}
	meta, err := s.Meta(ctx)
	if err != nil {
		return err
	}

	logrus.WithField("meta_value", meta).Debug("Read meta from workspace")
	metaTags, err := json.Marshal(meta.GetTags())
	if err != nil {
		return errors.Wrap(err, "could not marshal json")
	}
	err = workspace.SetVars("happymeta_", string(metaTags), "Happy Path metadata", false)
	if err != nil {
		return err
	}
	for k, v := range meta.GetParameters() {
		// TODO(el): why empty string?
		err = workspace.SetVars(k, v, "", false)
		if err != nil {
			return err
		}
	}
	workspace.ResetCache()

	tfDirPath := s.stackService.GetConfig().TerraformDirectory()

	happyProjectRoot := s.stackService.GetConfig().GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)

	logrus.Debugf("will use tf bundle found at %s", srcDir)

	tempFile, err := ioutil.TempFile("", "happy_tfe.*.tar.gz")
	if err != nil {
		return errors.Wrap(err, "could not create temporary file")
	}
	defer os.Remove(tempFile.Name())
	err = s.dirProcessor.Tarzip(srcDir, tempFile)
	if err != nil {
		return err
	}

	configVersionId, err := workspace.UploadVersion(srcDir, dryRun)
	if err != nil {
		return err
	}

	// TODO should be able to use workspace.Run() here, as workspace.UploadVersion(srcDir)
	// should have generated a Run containing the Config Version Id

	isDestroy := false
	err = workspace.RunConfigVersion(configVersionId, isDestroy, dryRun)
	if err != nil {
		return err
	}

	return workspace.WaitWithOptions(ctx, waitOptions, dryRun)
}

func (s *Stack) PrintOutputs(ctx context.Context) {
	logrus.Info("Module Outputs --")
	stackOutput, err := s.GetOutputs(ctx)
	if err != nil {
		logrus.Errorf("Failed to get output for stack %s: %s", s.stackName, err.Error())
		return
	}
	for k, v := range stackOutput {
		logrus.Printf("%s: %s", k, v)
	}
}

func (s *Stack) Print(ctx context.Context, name string, tablePrinter *util.TablePrinter) error {
	stackOutput, err := s.GetOutputs(ctx)
	if err != nil {
		return err
	}
	url := stackOutput["frontend_url"]
	status := s.GetStatus()
	meta, err := s.Meta(ctx)
	if err != nil {
		return err
	}

	tag := meta.DataMap["imagetag"]
	lastUpdated := meta.DataMap["updated"]
	imageTags, ok := meta.DataMap["imagetags"]
	if ok && len(imageTags) > 0 {
		var imageTagMap map[string]interface{}
		err = json.Unmarshal([]byte(imageTags), &imageTagMap)
		if err != nil {
			return errors.Wrap(err, "unable to parse json")
		}
		combinedTags := []string{tag}
		for imageTag := range imageTagMap {
			combinedTags = append(combinedTags, imageTag)
		}
		tag = strings.Join(combinedTags, ", ")
	}
	tablePrinter.AddRow(name, meta.DataMap["owner"], tag, status, url, lastUpdated)
	return nil
}
