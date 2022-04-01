package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Orchestrator struct {
	backend  *backend.Backend
	executor util.Executor
}

type container struct {
	host          string
	container     string
	arn           string
	taskID        string
	launchType    string
	containerName string
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

func (s *Orchestrator) Shell(ctx context.Context, stackName string, service string) error {
	clusterArn := s.backend.Conf().GetClusterArn()

	serviceName := stackName + "-" + service
	ecsClient := s.backend.GetECSClient()
	ec2Client := s.backend.GetEC2Client()

	listTaskInput := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterArn),
		ServiceName: aws.String(serviceName),
	}

	listTaskOutput, err := ecsClient.ListTasks(ctx, listTaskInput)
	if err != nil {
		return errors.Wrap(err, "error listing ecs tasks")
	}

	log.Println("Found tasks: ")
	headings := []string{"Task ID", "Started", "Status"}
	tablePrinter := util.NewTablePrinter(headings)

	describeTaskInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterArn),
		Tasks:   listTaskOutput.TaskArns,
	}

	describeTaskOutput, err := ecsClient.DescribeTasks(ctx, describeTaskInput)
	if err != nil {
		return errors.Wrap(err, "error describing ecs tasks")
	}

	containerMap := make(map[string]string)
	var containers []container

	for _, task := range describeTaskOutput.Tasks {
		taskArnSlice := strings.Split(*task.TaskArn, "/")
		taskID := taskArnSlice[len(taskArnSlice)-1]

		startedAt := "-"

		host := ""
		if task.ContainerInstanceArn != nil {
			host = *task.ContainerInstanceArn
		}

		if task.StartedAt != nil {
			startedAt = task.StartedAt.Format(time.RFC3339)
			containers = append(containers, container{
				host:          host,
				container:     *task.Containers[0].RuntimeId,
				arn:           *task.TaskArn,
				taskID:        taskID,
				launchType:    string(task.LaunchType),
				containerName: *task.Containers[0].Name,
			})
		}
		containerMap[*task.TaskArn] = host
		tablePrinter.AddRow(taskID, startedAt, *task.LastStatus)
	}

	tablePrinter.Print()
	// FIXME: we make the assumption of only one container in many places. need consistency
	// TODO: only support ECS exec-command and NOT SSH
	for _, container := range containers {
		if container.launchType == config.LaunchTypeFargate.String() {
			awsProfile := s.backend.Conf().AwsProfile()
			log.Infof("Connecting to %s:%s\n", container.taskID, container.containerName)
			// TODO: use the Go SDK and don't shell out
			//       see https://github.com/tedsmitt/ecsgo/blob/c1509097047a2d037577b128dcda4a35e23462fd/internal/pkg/internal.go#L196
			awsArgs := []string{"aws", "--profile", *awsProfile, "ecs", "execute-command", "--cluster", clusterArn, "--container", container.containerName, "--command", "/bin/bash", "--interactive", "--task", container.taskID}

			awsCmd, err := s.executor.LookPath("aws")
			if err != nil {
				return errors.Wrap(err, "failed to locate the AWS cli")
			}

			cmd := &exec.Cmd{
				Path:   awsCmd,
				Args:   awsArgs,
				Stdin:  os.Stdin,
				Stderr: os.Stderr,
				Stdout: os.Stdout,
			}
			log.Println(cmd)
			if err := s.executor.Run(cmd); err != nil {
				return errors.Wrap(err, "failed to execute")
			}
		}
		input := &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(clusterArn),
			ContainerInstances: []string{container.host},
		}

		result, err := ecsClient.DescribeContainerInstances(ctx, input)
		if err != nil {
			return errors.Wrap(err, "could not describe ecs container instances")
		}

		ec2InstanceId := result.ContainerInstances[0].Ec2InstanceId

		describeInstancesInput := &ec2.DescribeInstancesInput{
			InstanceIds: []string{*ec2InstanceId},
		}

		describeInstanceOutput, err := ec2Client.DescribeInstances(ctx, describeInstancesInput)
		if err != nil {
			return errors.Wrap(err, "could not describe instances")
		}

		ipAddress := describeInstanceOutput.Reservations[0].Instances[0].PrivateIpAddress

		log.Infof("Connecting to: %s %s\n", container.arn, *ipAddress)

		// FIXME: assumes /bin/bash present
		args := []string{
			"ssh", "-t", *ipAddress,
			"sudo", "docker", "exec", "-ti", container.container, "/bin/bash"}

		sshCmd, err := s.executor.LookPath("ssh")
		if err != nil {
			return errors.Wrap(err, "ssh not found in PATH")
		}

		cmd := &exec.Cmd{
			Path:   sshCmd,
			Args:   args,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		log.Infof("Command to connect: %s\n", cmd)
		//TODO: For now just print the commands to connect to
		// all the containers. Will make it a bit interactive
		// to select the container.
		// if err := cmd.Run(); err != nil {
		// 	return errors.Wrap(err, "Failed to ssh")
		// }
	}

	return nil
}

// Taking tasks defined in the config, look up their ID (e.g. ARN) in the given Stack
// object, and run these tasks with TaskRunner
func (s *Orchestrator) RunTasks(ctx context.Context, stack *stack_mgr.Stack, taskType string) error {
	taskOutputs, err := s.backend.Conf().GetTasks(taskType)
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

	for _, taskDef := range tasks {
		log.Infof("using task definition %s", taskDef)
		err = s.backend.RunTask(ctx, taskDef, launchType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Orchestrator) Logs(ctx context.Context, stackName string, serviceName string, since string, follow bool) error {
	endTime := time.Now()

	duration, err := time.ParseDuration(since)
	if err != nil {
		return errors.Wrapf(err, "invalid duration: '%s'", since)
	}
	startTime := endTime.Add(-duration)

	logGroup := ""
	logStreamName := ""
	serviceName = fmt.Sprintf("%s-%s", stackName, serviceName)

	taskArns, err := s.backend.GetTasks(ctx, &serviceName)
	if err != nil {
		return errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
	}

	for _, taskArn := range taskArns {
		resourceArn, err := arn.Parse(taskArn)
		if err != nil {
			return errors.Wrapf(err, "unavle to parse task ARN: '%s'", taskArn)
		}

		arnSegments := strings.Split(resourceArn.Resource, "/")
		if len(arnSegments) < 3 {
			continue
		}
		taskId := arnSegments[len(arnSegments)-1]

		taskDefinitions, err := s.backend.GetTaskDefinitions(ctx, taskArn)
		if err != nil {
			return errors.Wrapf(err, "error retrieving task definition for task '%s'", taskArn)
		}

		for _, taskDefinition := range taskDefinitions {
			for _, containerDefinition := range taskDefinition.ContainerDefinitions {
				logGroup = containerDefinition.LogConfiguration.Options["awslogs-group"]
				logStreamName = containerDefinition.LogConfiguration.Options["awslogs-stream-prefix"] + "/" + *containerDefinition.Name + "/" + taskId

				if len(logGroup) > 0 {
					break
				}
			}
			if len(logGroup) > 0 {
				break
			}
		}
	}

	if len(logGroup) == 0 {
		return errors.Errorf("unable to determine a log group for service '%s'", serviceName)
	}

	params := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStreamName,
		StartTime:     aws.Int64(startTime.UnixNano() / int64(time.Millisecond)),
		EndTime:       aws.Int64(endTime.UnixNano() / int64(time.Millisecond)),
	}

	if follow {
		log.Infof("Tailing logs: group=%s, stream=%s", logGroup, logStreamName)
	}

	resp, err := s.backend.GetLogEventsAPIClient().GetLogEvents(ctx, &params)
	if err != nil {
		return errors.Wrap(err, "cannot retrieve logs")
	}

	for _, event := range resp.Events {
		log.Info(*event.Message)
	}

	if !follow {
		return nil
	}

	for {
		params.NextToken = resp.NextBackwardToken
		resp, err = s.backend.GetLogEventsAPIClient().GetLogEvents(ctx, &params)
		if err != nil {
			return errors.Wrap(err, "cannot tail logs")
		}

		for _, event := range resp.Events {
			log.Info(*event.Message)
		}
		time.Sleep(10 * time.Second)
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

		log.Infof("Incomplete deployment of service %s / Current status %v:\n", *service.ServiceName, incomplete)

		deregistered := 0
		for index := range service.Events {
			event := service.Events[len(service.Events)-1-index]
			eventTime := event.CreatedAt
			if time.Since(*eventTime) < (time.Hour) {
				continue
			}

			message := regexp.MustCompile(`^\(service ([^ ]+)\)`).ReplaceAllString(*event.Message, "$1")
			message = regexp.MustCompile(`\(([^ ]+) .*?\)`).ReplaceAllString(message, "$1")
			message = regexp.MustCompile(`:.*`).ReplaceAllString(message, "$1")
			if strings.Contains(message, "deregistered") {
				deregistered++
			}

			log.Infof("  %s %s\n", eventTime.Format(time.RFC3339), message)
			if deregistered > 3 {
				log.Println()
				log.Println("Many \"deregistered\" events - please check to see whether your service is crashing:")
				serviceName := strings.Replace(*service.ServiceName, fmt.Sprintf("%s-", stack), "", 1)
				// FIXME: what is this reference to a local happy script?
				log.Infof("  ./scripts/happy --env %s logs %s %s", s.backend.Conf().GetEnv(), stack, serviceName)
			}
		}
	}
	return nil
}
