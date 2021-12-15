package orchestrator

import (
	// "bufio"

	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chanzuckerberg/happy-deploy/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chanzuckerberg/happy-deploy/pkg/backend"
	"github.com/chanzuckerberg/happy-deploy/pkg/config"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/happy-deploy/pkg/stack_mgr"
)

type Orchestrator struct {
	config     config.HappyConfigIface
	taskRunner backend.TaskRunner
}

type container struct {
	host      string
	container string
	arn       string
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

	fmt.Println("Found tasks: ")
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
		containers = append(containers, container{
			host:      *task.ContainerInstanceArn,
			container: *task.Containers[0].RuntimeId,
			arn:       *task.TaskArn,
		})

		taskArnSlice := strings.Split(*task.TaskArn, "/")
		taskID := taskArnSlice[len(taskArnSlice)-1]
		containerMap[*task.TaskArn] = *task.ContainerInstanceArn
		tablePrinter.AddRow([]string{taskID, task.StartedAt.Format("2006-01-02 15:04:05"), *task.LastStatus})
	}

	tablePrinter.AddRow([]string{"", "", ""})
	tablePrinter.Print()

	for _, container := range containers {
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

		fmt.Println("Connecting to:", container.arn, *ipAddress)

		args := []string{"ssh", "-t", *ipAddress, "sudo", "docker", "exec", "-ti", container.container, "/bin/bash"}

		sshCmd, _ := exec.LookPath("ssh")

		cmd := &exec.Cmd{
			Path:   sshCmd,
			Args:   args,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		fmt.Println("Command to connect:", cmd)
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
func (s *Orchestrator) RunTasks(stack *stack_mgr.Stack, taskType string, wait bool, showLogs bool) error {

	// taskOutputs, ok := s.config.GetData().Tasks[taskType]
	// if !ok {
	// 	return fmt.Errorf("Tasks of type %s not found", taskType)
	// }
	taskOutputs, err := s.config.GetTasks(taskType)
	if err != nil {
		return err
	}

	stackOutputs, err := stack.GetOutputs()
	if err != nil {
		return err
	}

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
		wait := true
		s.taskRunner.RunTask(taskDef, wait)
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

	awsCmd, _ := exec.LookPath("aws")

	cmd := &exec.Cmd{
		Path:   awsCmd,
		Args:   awsArgs,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	fmt.Println(cmd)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Failed to get logs from AWS:")
	}

	return nil
}
