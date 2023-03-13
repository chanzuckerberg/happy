package stack_mgr

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/diagnostics"
	workspacerepo "github.com/chanzuckerberg/happy/cli/pkg/workspace_repo"
	"github.com/chanzuckerberg/happy/shared/opts"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type StackServiceIface interface {
	NewStackMeta(stackName string) *StackMeta
	Add(ctx context.Context, stackName string, dryRun bool, options ...opts.RunOption) (*Stack, error)
	Remove(ctx context.Context, stackName string, dryRun bool, options ...opts.RunOption) error
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
	return fmt.Sprintf("/happy/%s/stacklist", s.backend.Conf().GetEnv())
}

func (s *StackService) GetNamespacedWritePath() string {
	return fmt.Sprintf("/happy/%s/%s/stacklist", s.backend.Conf().App(), s.backend.Conf().GetEnv())
}

func (s *StackService) WithBackend(backend *backend.Backend) *StackService {
	creatorWorkspaceName := fmt.Sprintf("env-%s", backend.Conf().GetEnv())

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
func (s *StackService) resync(ctx context.Context, wait bool, options ...opts.RunOption) error {
	log.Debug("resyncing new workspace...")
	log.Debugf("running creator workspace %s...", s.creatorWorkspaceName)
	creatorWorkspace, err := s.workspaceRepo.GetWorkspace(ctx, s.creatorWorkspaceName)
	if err != nil {
		return errors.Wrapf(err, "unable to get workspace %s", s.creatorWorkspaceName)
	}
	err = creatorWorkspace.Run(ctx, options...)
	if err != nil {
		return errors.Wrapf(err, "error running latest %s workspace version", s.creatorWorkspaceName)
	}
	if wait {
		return creatorWorkspace.Wait(ctx, options...)
	}
	return nil
}

func (s *StackService) Remove(ctx context.Context, stackName string, dryRun bool, options ...opts.RunOption) error {
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

	err = s.resync(ctx, false, options...)
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

	lockKey := s.GetNamespacedWritePath()
	lock, err := distributedLock.AcquireLock(ctx, lockKey)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.removeFromStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", lockKey))
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

	return s.writeStacklist(ctx, stackNamesList)
}

func (s *StackService) Add(ctx context.Context, stackName string, dryRun bool, options ...opts.RunOption) (*Stack, error) {
	if dryRun {
		log.Debugf("temporarily creating a TFE workspace for stack '%s'", stackName)
	} else {
		log.Debugf("creating stack '%s'", stackName)
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
		if err := s.resync(ctx, wait, options...); err != nil {
			return nil, err
		}
	}

	_, err = s.GetStackWorkspace(ctx, stackName)
	if err != nil {
		return nil, err
	}

	stack := s.getOrCreateStack(stackName)
	s.stacks[stackName] = stack

	return stack, nil
}

func (s *StackService) addToStacklistWithLock(ctx context.Context, stackName string) error {
	distributedLock, err := s.getDistributedLock()
	if err != nil {
		return err
	}
	defer distributedLock.Close(ctx)

	lockKey := s.GetNamespacedWritePath()
	lock, err := distributedLock.AcquireLock(ctx, lockKey)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.addToStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", lockKey))
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

	return s.writeStacklist(ctx, newStackNames)
}

func (s *StackService) writeStacklist(ctx context.Context, stackNames []string) error {
	sort.Strings(stackNames)

	stackNamesJson, err := json.Marshal(stackNames)
	if err != nil {
		return errors.Wrap(err, "unable to serialize stack list as json")
	}

	stackNamesStr := string(stackNamesJson)
	log.WithFields(log.Fields{"path": s.GetNamespacedWritePath(), "data": stackNamesStr}).Debug("Writing to paramstore...")
	if err := s.backend.ComputeBackend.WriteParam(ctx, s.GetNamespacedWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}
	log.WithFields(log.Fields{"path": s.GetWritePath(), "data": stackNamesStr}).Debug("Writing to paramstore...")
	if err := s.backend.ComputeBackend.WriteParam(ctx, s.GetWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}

	return nil
}

func (s *StackService) GetStacks(ctx context.Context) (map[string]*Stack, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetStacks")
	if s.stacks != nil {
		return s.stacks, nil
	}

	log.WithField("path", s.GetNamespacedWritePath()).Debug("Reading stacks from paramstore at path...")
	paramOutput, err := s.backend.ComputeBackend.GetParam(ctx, s.GetNamespacedWritePath())
	if err != nil && strings.Contains(err.Error(), "ParameterNotFound") {
		log.WithField("path", s.GetWritePath()).Debug("Reading stacks from paramstore at path...")
		paramOutput, err = s.backend.ComputeBackend.GetParam(ctx, s.GetWritePath())
	}
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
		s.stacks[stackName] = s.getOrCreateStack(stackName)
	}

	return s.stacks, nil
}

func (s *StackService) GetStack(ctx context.Context, stackName string) (*Stack, error) {
	existingStacks, err := s.GetStacks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stacks")
	}
	stack, ok := existingStacks[stackName]
	if !ok {
		return nil, errors.Errorf("stack %s doesn't exist", stackName)
	}

	return stack, nil
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

// TODO: getOrCreateStack -> GetOrCreate?
func (s *StackService) getOrCreateStack(stackName string) *Stack {
	if stack, ok := s.stacks[stackName]; ok {
		return stack
	}

	stack := &Stack{
		stackService: s,
		Name:         stackName,
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
