package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
)

type Orchestrator struct {
	config     config.HappyConfigIface
	taskRunner backend.TaskRunner
}

type container struct {
	host          string
	container     string
	arn           string
	taskID        string
	launchType    string
	containerName string
}

func NewOrchestrator(config config.HappyConfigIface, taskRunner backend.TaskRunner) *Orchestrator {
	return &Orchestrator{
		config:     config,
		taskRunner: taskRunner,
	}
}

func (s *Orchestrator) Shell(stackName string, service string) error {
	clusterArn, err := s.config.ClusterArn()
	if err != nil {
		return err
	}

	serviceName := stackName + "-" + service
	ecsClient := s.taskRunner.GetECSClient()
	ec2Client := s.taskRunner.GetEC2Client()

	listTaskInput := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterArn),
		ServiceName: aws.String(serviceName),
	}

	listTaskOutput, err := ecsClient.ListTasks(listTaskInput)
	if err != nil {
		return err
	}

	log.Println("Found tasks: ")
	headings := []string{"Task ID", "Started", "Status"}
	tablePrinter := util.NewTablePrinter(headings)

	describeTaskInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterArn),
		Tasks:   listTaskOutput.TaskArns,
	}

	describeTaskOutput, err := ecsClient.DescribeTasks(describeTaskInput)
	if err != nil {
		return err
	}

	containerMap := make(map[string]string)
	var containers []container

	for _, task := range describeTaskOutput.Tasks {
		taskArnSlice := strings.Split(*task.TaskArn, "/")
		taskID := taskArnSlice[len(taskArnSlice)-1]

		startedAt := "-"

		if task.StartedAt != nil {
			startedAt = task.StartedAt.Format(time.RFC3339)
			containers = append(containers, container{
				host:          *task.ContainerInstanceArn,
				container:     *task.Containers[0].RuntimeId,
				arn:           *task.TaskArn,
				taskID:        taskID,
				launchType:    *task.LaunchType,
				containerName: *task.Containers[0].Name,
			})
		}
		containerMap[*task.TaskArn] = *task.ContainerInstanceArn
		tablePrinter.AddRow([]string{taskID, startedAt, *task.LastStatus})

	}

	tablePrinter.AddRow([]string{"", "", ""})
	tablePrinter.Print()

	for _, container := range containers {
		if container.launchType == config.LaunchTypeFargate {
			awsProfile := s.config.AwsProfile()
			log.Printf("Connecting to %s:%s\n", container.taskID, container.containerName)
			awsArgs := []string{"aws", "--profile", awsProfile, "ecs", "execute-command", "--cluster", clusterArn, "--container", container.containerName, "--command", "/bin/bash", "--interactive", "--task", container.taskID}

			awsCmd, err := exec.LookPath("aws")
			if err != nil {
				return errors.Wrap(err, "failed to locate the AWS cli")
			}

			cmd := &exec.Cmd{
				Path:   awsCmd,
				Args:   awsArgs,
				Stderr: os.Stderr,
				Stdout: os.Stdout,
			}
			log.Println(cmd)
			if err := cmd.Run(); err != nil {
				return errors.Wrap(err, "failed to execute")
			}
		}
		input := &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(clusterArn),
			ContainerInstances: aws.StringSlice([]string{container.host}),
		}

		result, err := ecsClient.DescribeContainerInstances(input)
		if err != nil {
			return err
		}

		ec2InstanceId := result.ContainerInstances[0].Ec2InstanceId

		describeInstancesInput := &ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{*ec2InstanceId}),
		}

		describeInstanceOutput, err := ec2Client.DescribeInstances(describeInstancesInput)
		if err != nil {
			return err
		}

		ipAddress := describeInstanceOutput.Reservations[0].Instances[0].PrivateIpAddress

		log.Printf("Connecting to: %s %s\n", container.arn, *ipAddress)

		args := []string{"ssh", "-t", *ipAddress, "sudo", "docker", "exec", "-ti", container.container, "/bin/bash"}

		sshCmd, err := exec.LookPath("ssh")
		if err != nil {
			return err
		}

		cmd := &exec.Cmd{
			Path:   sshCmd,
			Args:   args,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		log.Printf("Command to connect: %s\n", cmd)
		//TODO: For now just print the commands to connect to
		// all the containers. Will make it a bit interactive
		// to select the container.
		// if err := cmd.Run(); err != nil {
		// 	return errors.Wrap(err, "Failed to ssh")
		// }
	}

	return nil
}

// Taking tasks defined in the config, look up their ID (e.g ARN) in the given Stack
// object, and run these tasks with TaskRunner
func (s *Orchestrator) RunTasks(stack *stack_mgr.Stack, taskType string, showLogs bool) error {
	taskOutputs, err := s.config.GetTasks(taskType)
	if err != nil {
		return err
	}

	stackOutputs, err := stack.GetOutputs()
	if err != nil {
		return err
	}

	launchType := s.config.TaskLaunchType()

	tasks := []string{}
	for _, taskOutput := range taskOutputs {
		task, ok := stackOutputs[taskOutput]
		if !ok {
			continue
		}
		tasks = append(tasks, task)
	}

	for _, taskDef := range tasks {
		fmt.Printf("Using task definition %s\n", taskDef)
		err = s.taskRunner.RunTask(taskDef, launchType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Orchestrator) Logs(stackName string, service string, since string) error {
	// TODO get logs path from ECS instead of generating
	logPrefix := s.config.LogGroupPrefix()
	logPath := fmt.Sprintf("%s/%s/%s", logPrefix, stackName, service)

	awsProfile := s.config.AwsProfile()
	regionName := "us-west-2"
	awsArgs := []string{"aws", "--profile", awsProfile, "--region", regionName, "logs", "tail", "--since", since, "--follow", logPath}

	awsCmd, err := exec.LookPath("aws")
	if err != nil {
		return errors.Wrap(err, "failed to locate the AWS cli")
	}

	cmd := &exec.Cmd{
		Path:   awsCmd,
		Args:   awsArgs,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	log.Println(cmd)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Failed to get logs from AWS:")
	}

	return nil
}
