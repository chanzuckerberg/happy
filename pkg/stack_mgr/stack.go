package stack_mgr

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type StackIface interface {
	GetName() string
	SetMeta(meta *StackMeta)
	Meta() (*StackMeta, error)
	GetStatus() string
	GetOutputs() (map[string]string, error)
	Wait(waitOptions options.WaitOptions)
	Apply(waitOptions options.WaitOptions) error
	Destroy() error
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

func (s *Stack) getWorkspace() (workspace_repo.Workspace, error) {
	if s.workspace == nil {
		workspace, err := s.stackService.GetStackWorkspace(s.stackName)
		if err != nil {
			return nil, errors.Errorf("failed to get workspace for stack %s", s.stackName)
		}
		s.workspace = workspace
	}

	return s.workspace, nil
}

func (s *Stack) GetOutputs() (map[string]string, error) {
	workspace, err := s.getWorkspace()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get output for stack %s", s.stackName)
	}

	outputs, err := workspace.GetOutputs()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get output for stack %s", s.stackName)
	}

	return outputs, nil
}

func (s *Stack) GetStatus() string {
	if s.workspace != nil {
		return s.workspace.GetCurrentRunStatus()
	}
	return ""
}

func (s *Stack) SetMeta(meta *StackMeta) {
	s.meta = meta
}

func (s *Stack) Meta() (*StackMeta, error) {
	if s.meta == nil {
		s.meta = s.stackService.NewStackMeta(s.stackName)

		// update tags of meta with those from the backing workspace
		workspace, err := s.getWorkspace()
		if err != nil {
			return nil, err
		}
		tags, err := workspace.GetTags()
		if err != nil {
			return nil, err
		}

		// TODO(el): what is this?
		if len(tags) == 0 {
			tags = map[string]string{
				"happy/meta/owner":    "UNKNOWN",
				"happy/meta/imagetag": "UNKNOWN",
			}
		}

		err = s.meta.Load(tags)
		if err != nil {
			return nil, err
		}
	}
	return s.meta, nil
}

func (s *Stack) Destroy() error {
	workspace, err := s.getWorkspace()
	if err != nil {
		return err
	}
	versionId, err := workspace.GetLatestConfigVersionID()

	// NOTE [hanxlin] I do not know, when last version does not exist, if the call
	// returns an error an an empty string. So it checks both here.
	if err != nil || versionId == "" {
		log.Warn("No latest version of workspace to destroy. Assuming already empty and continuing.")
		return nil
	}
	isDestroy := true
	err = workspace.Run(isDestroy)
	if err != nil {
		return err
	}

	return workspace.Wait()
}

func (s *Stack) Wait(waitOptions options.WaitOptions) error {
	workspace, err := s.getWorkspace()
	if err != nil {
		return err
	}
	return workspace.WaitWithOptions(waitOptions)
}

func (s *Stack) Apply(waitOptions options.WaitOptions) error {
	log.Infof("apply stack %s...", s.stackName)

	workspace, err := s.getWorkspace()
	if err != nil {
		return err
	}
	meta, err := s.Meta()
	if err != nil {
		return err
	}

	log.WithField("meta_value", meta).Debug("Read meta from workspace")
	metaTags, err := json.Marshal(meta.GetTags())
	if err != nil {
		return err
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

	// TODO: get this configuration earlier, by this point it's too late
	happyProjectRoot, _ := os.LookupEnv("HAPPY_PROJECT_ROOT")
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	tempFile, err := ioutil.TempFile(curDir, "happy_tfe.*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	err = s.dirProcessor.Tarzip(srcDir, tempFile)
	if err != nil {
		return err
	}

	configVersionId, err := workspace.UploadVersion(srcDir)
	if err != nil {
		return err
	}

	// TODO should be able to use workspace.Run() here, as workspace.UploadVersion(srcDir)
	// should have generated a Run containing the Config Version Id
	isDestroy := false
	err = workspace.RunConfigVersion(configVersionId, isDestroy)
	if err != nil {
		return err
	}

	return workspace.WaitWithOptions(waitOptions)
}

func (s *Stack) PrintOutputs() {
	logrus.Info("Module Outputs --")
	stackOutput, err := s.GetOutputs()
	if err != nil {
		logrus.Errorf("Failed to get output for stack %s: %s", s.stackName, err.Error())
	}
	for k, v := range stackOutput {
		logrus.Printf("%s: %s", k, v)
	}
}
