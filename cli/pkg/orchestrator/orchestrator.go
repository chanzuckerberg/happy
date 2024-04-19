package orchestrator

import (
	"context"
	"encoding/json"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Orchestrator struct {
	backend     *backend.Backend
	happyConfig *config.HappyConfig
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		happyConfig: nil,
	}
}

func (s *Orchestrator) WithBackend(backend *backend.Backend) *Orchestrator {
	s.backend = backend
	return s
}

func (s *Orchestrator) WithHappyConfig(happyConfig *config.HappyConfig) *Orchestrator {
	s.happyConfig = happyConfig
	return s
}

func (s *Orchestrator) Shell(ctx context.Context, stackName, serviceName, containerName string, shellCommand []string) error {
	return s.backend.Shell(ctx, stackName, serviceName, containerName, shellCommand)
}

func (s *Orchestrator) TaskExists(ctx context.Context, taskType backend.TaskType) bool {
	return s.happyConfig.TaskExists(string(taskType))
}

// Taking tasks defined in the config, look up their ID (e.g. ARN) in the given Stack
// object, and run these tasks with TaskRunner
func (s *Orchestrator) RunTasks(ctx context.Context, stack *stack.Stack, taskType backend.TaskType) error {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return nil
	}

	if !s.TaskExists(ctx, taskType) {
		log.Warnf("No tasks defined for type %s, skipping.", taskType)
		return nil
	}

	taskOutputs, err := s.happyConfig.GetTasks(string(taskType))
	if err != nil {
		return err
	}

	stackOutputs, err := stack.GetOutputs(ctx)
	if err != nil {
		return err
	}

	launchType := s.happyConfig.TaskLaunchType()

	tasks := []string{}
	if taskArns, ok := stackOutputs["task_arns"]; ok {
		var taskMap map[string]string
		err = json.Unmarshal([]byte(taskArns), &taskMap)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal task_arns: '%s'", taskArns)
		}
		for taskName, taskArn := range taskMap {
			for _, taskOutput := range taskOutputs {
				if taskName == taskOutput {
					tasks = append(tasks, taskArn)
				}
			}
		}
	} else {
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

func (s *Orchestrator) PrintLogs(ctx context.Context, stack string, services []string) error {
	for _, service := range services {
		log.Infof("Printing logs for service %s", service)
		err := s.backend.PrintLogs(ctx, stack, service, "")
		if err != nil {
			log.Errorf("Failed to print logs for service %s: %s\n", service, err.Error())
		}
	}
	return nil
}

func (s *Orchestrator) GetResources(ctx context.Context, stack *stack.Stack) ([]util.ManagedResource, error) {
	resources, err := stack.GetResources(ctx)
	if err != nil {
		return make([]util.ManagedResource, 0), err
	}
	return resources, nil
}
