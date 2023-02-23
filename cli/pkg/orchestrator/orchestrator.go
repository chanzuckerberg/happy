package orchestrator

import (
	"context"

	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	log "github.com/sirupsen/logrus"
)

type Orchestrator struct {
	backend *backend.Backend
	dryRun  bool
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (s *Orchestrator) WithBackend(backend *backend.Backend) *Orchestrator {
	s.backend = backend
	return s
}

func (s *Orchestrator) WithDryRun(dryRun bool) *Orchestrator {
	s.dryRun = dryRun
	return s
}

func (s *Orchestrator) Shell(ctx context.Context, stackName string, service string) error {
	return s.backend.Shell(ctx, stackName, service)
}

func (s *Orchestrator) TaskExists(ctx context.Context, taskType backend.TaskType) bool {
	return s.backend.Conf().TaskExists(string(taskType))
}

// Taking tasks defined in the config, look up their ID (e.g. ARN) in the given Stack
// object, and run these tasks with TaskRunner
func (s *Orchestrator) RunTasks(ctx context.Context, stack *stack_mgr.Stack, taskType backend.TaskType) error {
	if s.dryRun {
		return nil
	}

	if !s.TaskExists(ctx, taskType) {
		log.Warnf("No tasks defined for type %s, skipping.", taskType)
		return nil
	}

	taskOutputs, err := s.backend.Conf().GetTasks(string(taskType))
	if err != nil {
		return err
	}

	stackOutputs, err := stack.GetOutputs(ctx)
	if err != nil {
		return err
	}

	launchType := s.backend.Conf().TaskLaunchType()

	tasks := []string{}
	for _, taskOutput := range taskOutputs {
		task, ok := stackOutputs[taskOutput]
		if !ok {
			continue
		}
		if len(task) >= 2 {
			if task[0] == '"' && task[len(task)-1] == '"' {
				task = task[1 : len(task)-1]
			}
		}
		tasks = append(tasks, task)
	}

	log.Infof("running after update tasks %+v", tasks)
	for _, taskDef := range tasks {
		err = s.backend.RunTask(ctx, taskDef, launchType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Orchestrator) GetEvents(ctx context.Context, stack string, services []string) error {
	return s.backend.GetEvents(ctx, stack, services)
}
