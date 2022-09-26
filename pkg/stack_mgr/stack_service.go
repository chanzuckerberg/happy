package stack_mgr

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/util"
	workspacerepo "github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type StackServiceIface interface {
	NewStackMeta(stackName string) *StackMeta
	Add(ctx context.Context, stackName string, dryRun util.DryRunType) (*Stack, error)
	Remove(ctx context.Context, stackName string, dryRun util.DryRunType) error
	GetStacks(ctx context.Context) (map[string]*Stack, error)
	GetStackWorkspace(ctx context.Context, stackName string) (workspacerepo.Workspace, error)
	GetConfig() *config.HappyConfig
}

type StackService struct {
	// dependencies
	backend       *backend.Backend
	workspaceRepo workspacerepo.WorkspaceRepoIface
	dirProcessor  util.DirProcessor
	executor      util.Executor

	// attributes
	writePath string

	// NOTE: creator Workspace is a workspace that creates dependent workspaces with
	// given default values and configuration
	// the derived workspace is then used to launch the actual happy infrastructure
	creatorWorkspaceName string

	// cache
	stacks map[string]*Stack
}

func NewStackService() *StackService {
	// TODO pass this in instead?
	dirProcessor := util.NewLocalProcessor()

	return &StackService{
		stacks:       nil,
		dirProcessor: dirProcessor,
		executor:     util.NewDefaultExecutor(),
	}
}

func (s *StackService) GetWritePath() string {
	return s.writePath
}

func (s *StackService) WithBackend(backend *backend.Backend) *StackService {
	creatorWorkspaceName := fmt.Sprintf("env-%s", backend.Conf().GetEnv())

	s.writePath = backend.Conf().GetSsmStacklistParamPath()
	if s.writePath == "" {
		// use the default value if no custom path is set
		s.writePath = fmt.Sprintf("/happy/%s/stacklist", backend.Conf().GetEnv())
	}
	s.creatorWorkspaceName = creatorWorkspaceName
	s.backend = backend

	return s
}

func (s *StackService) WithExecutor(executor util.Executor) *StackService {
	s.executor = executor
	return s
}

func (s *StackService) WithWorkspaceRepo(workspaceRepo workspacerepo.WorkspaceRepoIface) *StackService {
	s.workspaceRepo = workspaceRepo
	return s
}

func (s *StackService) NewStackMeta(stackName string) *StackMeta {
	// TODO: what are all these translations?
	dataMap := map[string]string{
		"app":      s.backend.Conf().App(),
		"env":      s.backend.Conf().GetEnv(),
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
		StackName: stackName,
		DataMap:   dataMap,
		TagMap:    tagMap,
		ParamMap:  paramMap,
	}
}

func (s *StackService) GetConfig() *config.HappyConfig {
	return &s.backend.Conf().HappyConfig
}

// Invoke a specific TFE workspace that creates/deletes TFE workspaces,
// with prepopulated variables for identifier tokens.
func (s *StackService) resync(ctx context.Context, wait bool) error {
	log.Debug("resyncing new workspace...")
	log.Debugf("running creator workspace %s...", s.creatorWorkspaceName)
	creatorWorkspace, err := s.workspaceRepo.GetWorkspace(ctx, s.creatorWorkspaceName)
	if err != nil {
		return err
	}
	isDestroy := false
	err = creatorWorkspace.Run(isDestroy, false)
	if err != nil {
		return err
	}
	if wait {
		return creatorWorkspace.Wait(ctx, false)
	}
	return nil
}

func (s *StackService) Remove(ctx context.Context, stackName string, dryRun util.DryRunType) error {
	if dryRun {
		return nil
	}
	var err error
	if s.GetConfig().GetFeatures().EnableDynamoLocking {
		err = s.removeFromStacklistWithLock(ctx, stackName)
	} else {
		err = s.removeFromStacklist(ctx, stackName)
	}
	if err != nil {
		return err
	}

	wait := false // no need to wait for TFE workspace to finish removing
	err = s.resync(ctx, wait)
	if err != nil {
		return errors.Wrap(err, "unable to resync the workspace")
	}
	delete(s.stacks, stackName)

	return nil
}

func (s *StackService) removeFromStacklistWithLock(ctx context.Context, stackName string) error {
	distributedLock, err := s.getDistributedLock()
	if err != nil {
		return err
	}
	defer distributedLock.Close(ctx)

	lock, err := distributedLock.AcquireLock(ctx, s.writePath)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.removeFromStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", s.writePath))
	}

	return ret
}

func (s *StackService) removeFromStacklist(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Removing stack...")

	s.stacks = nil // force a refresh of stacks.
	stacks, err := s.GetStacks(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get a list of stacks")
	}
	stackNamesList := []string{}
	for name := range stacks {
		if name != stackName {
			stackNamesList = append(stackNamesList, name)
		}
	}

	sort.Strings(stackNamesList)
	stackNamesJson, err := json.Marshal(stackNamesList)
	if err != nil {
		return errors.Wrap(err, "unable to serialize stack list as json")
	}
	err = s.backend.ComputeBackend.WriteParam(ctx, s.writePath, string(stackNamesJson))
	if err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}

	return nil
}

func (s *StackService) Add(ctx context.Context, stackName string, dryRun util.DryRunType) (*Stack, error) {
	if dryRun {
		log.Infof("temporarily creating a TFE workspace for stack '%s'", stackName)
	} else {
		log.Infof("creating stack '%s'", stackName)
	}

	var err error
	if s.GetConfig().GetFeatures().EnableDynamoLocking {
		err = s.addToStacklistWithLock(ctx, stackName)
	} else {
		err = s.addToStacklist(ctx, stackName)
	}
	if err != nil {
		return nil, err
	}

	if !util.IsLocalstackMode() {
		// Create the workspace
		wait := true
		if err := s.resync(ctx, wait); err != nil {
			return nil, err
		}
	}

	if _, err := s.GetStackWorkspace(ctx, stackName); err != nil {
		return nil, err
	}

	stack := s.GetStack(stackName)
	s.stacks[stackName] = stack

	return stack, nil
}

func (s *StackService) addToStacklistWithLock(ctx context.Context, stackName string) error {
	distributedLock, err := s.getDistributedLock()
	if err != nil {
		return err
	}
	defer distributedLock.Close(ctx)

	lock, err := distributedLock.AcquireLock(ctx, s.writePath)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.addToStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", s.writePath))
	}

	return ret
}

func (s *StackService) addToStacklist(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Adding new stack...")

	// force refresh list of stacks, and add to it the new stack
	s.stacks = nil
	existStacks, err := s.GetStacks(ctx)
	if err != nil {
		return err
	}

	newStackNames := []string{}
	stackNameExists := false
	for name := range existStacks {
		newStackNames = append(newStackNames, name)
		if name == stackName {
			stackNameExists = true
		}
	}
	if !stackNameExists {
		newStackNames = append(newStackNames, stackName)
	}

	sort.Strings(newStackNames)

	stackNamesJson, err := json.Marshal(newStackNames)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"path": s.writePath,
		"data": stackNamesJson,
	}).Debug("Writing to paramstore...")
	if err := s.backend.ComputeBackend.WriteParam(ctx, s.writePath, string(stackNamesJson)); err != nil {
		return err
	}

	return nil
}

func (s *StackService) GetStacks(ctx context.Context) (map[string]*Stack, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetStacks")
	if s.stacks != nil {
		return s.stacks, nil
	}

	log.WithField("path", s.writePath).Debug("Reading stacks from paramstore at path...")
	paramOutput, err := s.backend.ComputeBackend.GetParam(ctx, s.writePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stacks")
	}

	log.WithField("output", paramOutput).Debug("read stacks info from param store")

	var stacklist []string
	err = json.Unmarshal([]byte(paramOutput), &stacklist)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	log.WithField("output", stacklist).Debug("marshalled json output to string slice")

	s.stacks = map[string]*Stack{}
	for _, stackName := range stacklist {
		s.stacks[stackName] = s.GetStack(stackName)
	}

	return s.stacks, nil
}

// pre-format stack name and call workspaceRepo's GetWorkspace method
func (s *StackService) GetStackWorkspace(ctx context.Context, stackName string) (workspacerepo.Workspace, error) {
	workspaceName := fmt.Sprintf("%s-%s", s.backend.Conf().GetEnv(), stackName)

	ws, err := s.workspaceRepo.GetWorkspace(ctx, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workspace")
	}

	return ws, nil
}

// TODO: GetStack -> GetOrCreate?
func (s *StackService) GetStack(stackName string) *Stack {
	if stack, ok := s.stacks[stackName]; ok {
		return stack
	}

	stack := &Stack{
		stackService: s,
		stackName:    stackName,
		dirProcessor: s.dirProcessor,
		executor:     s.executor,
	}

	s.stacks[stackName] = stack

	return stack
}

func (s *StackService) HasState(ctx context.Context, stackName string) (bool, error) {
	workspace, err := s.GetStackWorkspace(ctx, stackName)
	if err != nil {
		if errors.Is(err, tfe.ErrInvalidWorkspaceValue) || errors.Is(err, tfe.ErrResourceNotFound) {
			// Workspace doesn't exist, thus no state
			return false, nil
		}
		return true, errors.Wrap(err, "Cannot get the stack workspace")
	}
	return workspace.HasState(ctx)
}

func (s *StackService) getDistributedLock() (*backend.DistributedLock, error) {
	lockConfig := backend.DistributedLockConfig{DynamodbTableName: s.backend.Conf().GetDynamoLocktableName()}
	return backend.NewDistributedLock(&lockConfig, s.backend.GetDynamoDBClient())
}
