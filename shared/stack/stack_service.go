package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	workspacerepo "github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
)

type StackServiceIface interface {
	NewStackMeta(stackName string) *StackMeta
	Add(ctx context.Context, stackName string, options ...workspacerepo.TFERunOption) (*Stack, error)
	Remove(ctx context.Context, stackName string, options ...workspacerepo.TFERunOption) error
	GetStacks(ctx context.Context) (map[string]*Stack, error)
	GetStackWorkspace(ctx context.Context, stackName string) (workspacerepo.Workspace, error)
}

type StackService struct {
	backend       *backend.Backend
	workspaceRepo workspacerepo.WorkspaceRepoIface
	executor      util.Executor
	env, appName  string

	// NOTE: creator Workspace is a workspace that creates dependent workspaces with
	// given default values and configuration
	// the derived workspace is then used to launch the actual happy infrastructure
	creatorWorkspaceName string
}

func NewStackService(env, appName string) *StackService {
	return &StackService{
		executor: util.NewDefaultExecutor(),
		env:      env,
		appName:  appName,
	}
}

func (s *StackService) GetWritePath() string {
	return fmt.Sprintf("/happy/%s/stacklist", s.env)
}

func (s *StackService) GetNamespacedWritePath() string {
	return fmt.Sprintf("/happy/%s/%s/stacklist", s.appName, s.env)
}

func (s *StackService) WithBackend(backend *backend.Backend) *StackService {
	creatorWorkspaceName := fmt.Sprintf("env-%s", s.env)

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

// Invoke a specific TFE workspace that creates/deletes TFE workspaces,
// with prepopulated variables for identifier tokens.
func (s *StackService) resync(ctx context.Context, options ...workspacerepo.TFERunOption) error {
	log.Debug("resyncing new workspace...")
	log.Debugf("running creator workspace %s...", s.creatorWorkspaceName)
	creatorWorkspace, err := s.workspaceRepo.GetWorkspace(ctx, s.creatorWorkspaceName)
	if err != nil {
		return errors.Wrapf(err, "unable to get workspace %s", s.creatorWorkspaceName)
	}
	err = creatorWorkspace.Run(ctx, options...)
	if err != nil {
		return errors.Wrapf(err, "error running latest %s workspace version, examine the plan: %s", s.creatorWorkspaceName, creatorWorkspace.GetWorkspaceUrl())
	}
	err = creatorWorkspace.Wait(ctx)
	if err != nil {
		return errors.Wrapf(err, "error waiting on the conclusion of latest %s workspace plan, please examine the output: %s", s.creatorWorkspaceName, creatorWorkspace.GetWorkspaceUrl())
	}
	return nil
}

func (s *StackService) GetLatestDeployedTag(ctx context.Context, stackName string) (string, error) {
	stack, err := s.GetStack(ctx, stackName)
	if err != nil {
		return "", errors.Wrap(err, "unable to get the stack")
	}
	stackInfo, err := stack.GetStackInfo(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get the stack info")
	}
	return stackInfo.Tag, nil
}

func (s *StackService) Remove(ctx context.Context, stackName string, opts ...workspacerepo.TFERunOption) error {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return nil
	}
	enableDynamoLocking, ok := ctx.Value(options.EnableDynamoLockingKey).(bool)
	if !ok {
		// default to true, with the option to override.
		// this is an old enough feature where I don't think we need to have behind a feature flag
		// all happy environments come with a dynamo table
		enableDynamoLocking = true
	}
	var err error
	if enableDynamoLocking {
		err = s.removeFromStacklistWithLock(ctx, stackName)
	} else {
		err = s.removeFromStacklist(ctx, stackName)
	}
	if err != nil {
		return errors.Wrap(err, "unable to remove stack from stacklist")
	}

	err = s.resync(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "Removal of the stack workspace failed, but stack was removed from the stack list. Please examine the plan")
	}
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

func (s *StackService) Add(ctx context.Context, stackName string, opts ...workspacerepo.TFERunOption) (*Stack, error) {
	log.WithField("stack_name", stackName).Debug("Adding a new stack...")
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		log.Debugf("temporarily creating a TFE workspace for stack '%s'", stackName)
	} else {
		log.Debugf("creating stack '%s'", stackName)
	}

	enableDynamoLocking, ok := ctx.Value(options.EnableDynamoLockingKey).(bool)
	if !ok {
		// default to true, with the option to override.
		// this is an old enough feature where I don't think we need to have behind a feature flag
		// all happy environments come with a dynamo table
		enableDynamoLocking = true
	}
	var err error
	if enableDynamoLocking {
		err = s.addToStacklistWithLock(ctx, stackName)
	} else {
		err = s.addToStacklist(ctx, stackName)
	}
	if err != nil {
		return nil, err
	}

	if !util.IsLocalstackMode() {
		// Create the workspace
		if err := s.resync(ctx, opts...); err != nil {
			return nil, err
		}
	}

	_, err = s.GetStackWorkspace(ctx, stackName)
	if err != nil {
		return nil, err
	}
	return s.createStack(stackName), nil
}

func (s *StackService) addToStacklistWithLock(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Adding new stack with a lock...")
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
	cb, err := s.backend.GetComputeBackend(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to connect to a compute backend")
	}

	if err := cb.WriteParam(ctx, s.GetNamespacedWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}
	log.WithFields(log.Fields{"path": s.GetWritePath(), "data": stackNamesStr}).Debug("Writing to paramstore...")
	if err := cb.WriteParam(ctx, s.GetWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}

	return nil
}

func (s *StackService) GetStacks(ctx context.Context) (map[string]*Stack, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetStacks")
	log.WithField("path", s.GetNamespacedWritePath()).Debug("Reading stacks from paramstore at path...")

	cb, err := s.backend.GetComputeBackend(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to a compute backend")
	}

	paramOutput, err := cb.GetParam(ctx, s.GetNamespacedWritePath())
	if err != nil && strings.Contains(err.Error(), "ParameterNotFound") {
		log.WithField("path", s.GetWritePath()).Debug("Reading stacks from paramstore at path...")
		paramOutput, err = cb.GetParam(ctx, s.GetWritePath())
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

	stacks := map[string]*Stack{}
	for _, stackName := range stacklist {
		stacks[stackName] = s.createStack(stackName)
	}

	return stacks, nil
}

func (s *StackService) CollectStackInfo(ctx context.Context, app string, listAll bool) ([]*model.AppStackResponse, error) {
	stacks, err := s.GetStacks(ctx)
	if err != nil {
		return nil, err
	}
	// Iterate in order
	stackNames := maps.Keys(stacks)
	stackInfos := make([]*model.AppStackResponse, len(stackNames))
	sort.Strings(stackNames)
	g, ctx := errgroup.WithContext(ctx)
	for i, name := range stackNames {
		i, name := i, name // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			stackInfo, err := stacks[name].GetStackInfo(ctx)
			stackInfos[i] = stackInfo
			if err != nil {
				log.Warnf("unable to get stack info for %s: %s (likely means the deploy failed the first time)", name, err)
				stackInfos[i] = &model.AppStackResponse{
					AppMetadata: *model.NewAppMetadata(app, s.env, name),
					StackMetadata: model.StackMetadata{
						TFEWorkspaceStatus: "error",
						Message:            err.Error(),
					},
					Error: err.Error(),
				}
			}

			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stack infos")
	}

	results := []*model.AppStackResponse{}
	for _, stackInfo := range stackInfos {
		if stackInfo == nil {
			// remove empties
			continue
		}
		// only show the stacks that belong to this app or they want to list all
		if listAll || (stackInfo.AppMetadata.App.AppName == app) {
			results = append(results, stackInfo)
		}
	}
	return results, nil
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
	workspaceName := fmt.Sprintf("%s-%s", s.env, stackName)

	ws, err := s.workspaceRepo.GetWorkspace(ctx, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workspace")
	}

	return ws, nil
}

func (s *StackService) createStack(stackName string) *Stack {
	return &Stack{
		stackService: s,
		Name:         stackName,
		executor:     s.executor,
	}
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
