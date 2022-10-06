package orchestrator

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Orchestrator struct {
	backend  *backend.Backend
	executor util.Executor
	dryRun   util.DryRunType
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		executor: util.NewDefaultExecutor(),
	}
}

func (s *Orchestrator) WithBackend(backend *backend.Backend) *Orchestrator {
	s.backend = backend
	return s
}

func (s *Orchestrator) WithExecutor(executor util.Executor) *Orchestrator {
	s.executor = executor
	return s
}

func (s *Orchestrator) WithDryRun(dryRun util.DryRunType) *Orchestrator {
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
		tasks = append(tasks, task)
	}

	log.Infof("running after update tasks %+v", tasks)
	for _, taskDef := range tasks {
		// TODO This needs to be pushed into the compute backend
		err = s.backend.RunTask(ctx, taskDef, launchType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Orchestrator) GetEvents(ctx context.Context, stack string, services []string) error {
	if len(services) == 0 {
		return nil
	}

	clusterArn := s.backend.Conf().GetClusterArn()

	ecsClient := s.backend.GetECSClient()

	ecsServices := make([]string, 0)
	for _, service := range services {
		ecsService := fmt.Sprintf("%s-%s", stack, service)
		ecsServices = append(ecsServices, ecsService)
	}

	describeServicesInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterArn),
		Services: ecsServices,
	}

	describeServicesOutput, err := ecsClient.DescribeServices(ctx, describeServicesInput)
	if err != nil {
		return errors.Wrap(err, "cannot describe services:")
	}

	for _, service := range describeServicesOutput.Services {
		incomplete := make([]string, 0)
		for _, deploy := range service.Deployments {
			if deploy.RolloutState != ecstypes.DeploymentRolloutStateCompleted {
				incomplete = append(incomplete, string(deploy.RolloutState))
			}
		}
		if len(incomplete) == 0 {
			continue
		}

		log.Println()
		log.Infof("Incomplete deployment of service %s / Current status %v:", *service.ServiceName, incomplete)

		deregistered := 0
		for index := range service.Events {
			event := service.Events[len(service.Events)-1-index]
			eventTime := event.CreatedAt
			if time.Since(*eventTime) > (time.Hour) {
				continue
			}

			message := regexp.MustCompile(`^\(service ([^ ]+)\)`).ReplaceAllString(*event.Message, "$1")
			message = regexp.MustCompile(`\(([^ ]+) .*?\)`).ReplaceAllString(message, "$1")
			message = regexp.MustCompile(`:.*`).ReplaceAllString(message, "$1")
			if strings.Contains(message, "deregistered") {
				deregistered++
			}

			log.Infof("  %s %s", eventTime.Format(time.RFC3339), message)
		}
		if deregistered > 3 {
			log.Println()
			log.Println("Many \"deregistered\" events - please check to see whether your service is crashing:")
			serviceName := strings.Replace(*service.ServiceName, fmt.Sprintf("%s-", stack), "", 1)
			log.Infof("  happy --env %s logs %s %s", s.backend.Conf().GetEnv(), stack, serviceName)
		}
	}
	return nil
}
