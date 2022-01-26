package stack_mgr

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	workspace_repo "github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type StackServiceIface interface {
	NewStackMeta(stackName string) *StackMeta
	Add(stackName string) (*Stack, error)
	Remove(stack_name string) error
	GetStacks() (map[string]*Stack, error)
	GetStackWorkspace(stackName string) (workspace_repo.Workspace, error)
	GetConfig() config.HappyConfigIface
}

type StackService struct {
	// dependencies
	config        config.HappyConfigIface
	paramStore    backend.ParamStoreBackend
	workspaceRepo workspace_repo.WorkspaceRepoIface
	dirProcessor  util.DirProcessor

	// attributes
	writePath            string
	creatorWorkspaceName string

	// cache
	stacks map[string]*Stack
}

func NewStackService(config config.HappyConfigIface, paramStore backend.ParamStoreBackend, workspaceRepo workspace_repo.WorkspaceRepoIface) *StackService {

	// TODO pass this in instead?
	dirProcessor := util.NewLocalProcessor()

	writePath := fmt.Sprintf("/happy/%s/stacklist", config.DefaultEnv())
	creatorWorkspaceName := fmt.Sprintf("env-%s", config.DefaultEnv())

	return &StackService{
		config:               config,
		writePath:            writePath,
		stacks:               nil,
		creatorWorkspaceName: creatorWorkspaceName,
		workspaceRepo:        workspaceRepo,
		paramStore:           paramStore,
		dirProcessor:         dirProcessor,
	}
}

func (s *StackService) NewStackMeta(stackName string) *StackMeta {
	dataMap := map[string]string{
		"app":      s.config.App(),
		"env":      s.config.DefaultEnv(),
		"instance": stackName,
	}

	tagMap := map[string]string{
		"app":          "happy/app",
		"env":          "happy/env",
		"instance":     "happy/instance",
		"owner":        "happy/meta/owner",
		"priority":     "happy/meta/priority",
		"slice":        "happy/meta/slice",
		"imagetag":     "happy/meta/imagetag",
		"imagetags":    "happy/meta/imagetags",
		"configsecret": "happy/meta/configsecret",
		"created":      "happy/meta/created-at",
		"updated":      "happy/meta/updated-at",
	}

	paramMap := map[string]string{
		"instance":     "stack_name",
		"slice":        "slice",
		"priority":     "priority",
		"imagetag":     "image_tag",
		"imagetags":    "image_tags",
		"configsecret": "happy_config_secret",
	}

	return &StackMeta{
		stackName: stackName,
		DataMap:   dataMap,
		TagMap:    tagMap,
		paramMap:  paramMap,
	}
}

func (s *StackService) GetConfig() config.HappyConfigIface {
	return s.config
}

// Invoke a specific TFE workspace that creates/deletes TFE workspaces,
// with prepopulated variables for identifier tokens.
func (s *StackService) resync(wait bool) error {
	log.Debug("Resyncing new workspace...")

	log.WithField("workspace_name", s.creatorWorkspaceName).Debug("Running workspace...")
	creatorWorkspace, err := s.workspaceRepo.GetWorkspace(s.creatorWorkspaceName)
	if err != nil {
		return err
	}
	isDestroy := false
	err = creatorWorkspace.Run(isDestroy)
	if err != nil {
		return err
	}
	if wait {
		return creatorWorkspace.Wait()
	}
	return nil
}

func (s *StackService) Remove(stackName string) error {
	log.WithField("stack_name", stackName).Debug("Removing stack...")

	s.stacks = nil // force a refresh of stacks.
	stacks, err := s.GetStacks()
	if err != nil {
		return err
	}
	stackNames := map[string]bool{}
	for currentStackName := range stacks {
		if currentStackName != stackName {
			stackNames[currentStackName] = true
		}
	}
	stackNamesList := []string{}
	for stackName := range stackNames {
		stackNamesList = append(stackNamesList, stackName)
	}

	sort.Strings(stackNamesList)
	stackNamesJson, err := json.Marshal(stackNamesList)
	if err != nil {
		return err
	}
	if err := s.paramStore.AddParams(s.writePath, string(stackNamesJson)); err != nil {
		return err
	}

	wait := false // no need to wait for TFE workspace to finish removing
	s.resync(wait)
	delete(stacks, stackName)

	return nil
}

func (s *StackService) Add(stackName string) (*Stack, error) {
	log.WithField("stack_name", stackName).Debug("Adding new stack...")

	// force refresh list of stacks, and add to it the new stack
	s.stacks = nil
	existStackNames := map[string]bool{}
	existStacks, err := s.GetStacks()
	if err != nil {
		return nil, err
	}
	for name := range existStacks {
		existStackNames[name] = true
	}
	existStackNames[stackName] = true
	newStackNames := []string{}
	for name := range existStackNames {
		newStackNames = append(newStackNames, name)
	}
	sort.Strings(newStackNames)

	stackNamesJson, err := json.Marshal(newStackNames)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"path": s.writePath,
		"data": stackNamesJson,
	}).Debug("Writing to paramstore...")
	if err := s.paramStore.AddParams(s.writePath, string(stackNamesJson)); err != nil {
		return nil, err
	}

	// Create the workspace
	wait := true
	if err := s.resync(wait); err != nil {
		return nil, err
	}

	if _, err := s.GetStackWorkspace(stackName); err != nil {
		return nil, err
	}

	stack := s.GetStack(stackName)
	s.stacks[stackName] = stack

	return stack, nil
}

func (s *StackService) GetStacks() (map[string]*Stack, error) {

	if s.stacks != nil {
		return s.stacks, nil
	}

	log.WithField("path", s.writePath).Debug("Reading stacks from paramstore at path...")
	paramOutput, err := s.paramStore.GetParameter(s.writePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get stacks")
	}

	log.WithField("output", *paramOutput).Debug("Read stacks info from param store")

	var stacklist []string
	json.Unmarshal([]byte(*paramOutput), &stacklist)

	log.WithField("output", stacklist).Debug("Marshalled json output to string slice")

	s.stacks = map[string]*Stack{}
	for _, stackName := range stacklist {
		s.stacks[stackName] = s.GetStack(stackName)
	}

	return s.stacks, nil
}

// pre-format stack name and call workspaceRepo's GetWorkspace method
func (s *StackService) GetStackWorkspace(stackName string) (workspace_repo.Workspace, error) {
	// TODO: check if env is passed to cmd
	workspaceName := fmt.Sprintf("%s-%s", s.config.DefaultEnv(), stackName)

	ws, err := s.workspaceRepo.GetWorkspace(workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to get workspace")
	}
	return ws, nil
}

func (s *StackService) GetStack(stackName string) *Stack {
	if stack, ok := s.stacks[stackName]; ok {
		return stack
	}

	stack := &Stack{
		stackService: s,
		stackName:    stackName,
		dirProcessor: s.dirProcessor,
	}

	s.stacks[stackName] = stack

	return stack
}
