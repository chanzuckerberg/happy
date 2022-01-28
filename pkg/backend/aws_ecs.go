package backend

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/chanzuckerberg/happy/pkg/config"
)

type TaskType string

const (
	DeletionTask  TaskType = "delete"
	MigrationTask TaskType = "migrate"
)

type TaskRunner interface {
	RunTask(taskDef string, launchType string, wait bool) error
	GetECSClient() ecsiface.ECSAPI
	GetEC2Client() *ec2.EC2
}

type AwsEcs struct {
	session   *session.Session
	ecsClient ecsiface.ECSAPI
	awsConfig *aws.Config
	config    config.HappyConfigIface
	logsSrc   *Cloudwatchlags
	ec2Client *ec2.EC2
}

var ecsSessInst TaskRunner
var creatECSMOnce sync.Once

func (s *AwsEcs) GetECSClient() ecsiface.ECSAPI {
	return s.ecsClient
}

func (s *AwsEcs) GetEC2Client() *ec2.EC2 {
	return s.ec2Client
}

func GetAwsEcs(config config.HappyConfigIface) TaskRunner {
	awsProfile := config.AwsProfile()
	creatECSMOnce.Do(func() {
		awsConfig := &aws.Config{
			Region:     aws.String("us-west-2"),
			MaxRetries: aws.Int(2),
		}
		session := session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           awsProfile,
			Config:            *awsConfig,
			SharedConfigState: session.SharedConfigEnable,
		}))

		ecsClient := ecs.New(session)
		ec2Client := ec2.New(session)

		logsSrc := GetAwsLogs(config)
		ecsSessInst = &AwsEcs{
			session:   session,
			ecsClient: ecsClient,
			awsConfig: awsConfig,
			config:    config,
			logsSrc:   logsSrc,
			ec2Client: ec2Client,
		}
	})

	return ecsSessInst
}

func (s *AwsEcs) RunTask(taskDefArn string, launchType string, wait bool) error {

	fmt.Printf("Running tasks for %s\n", taskDefArn)

	clusterArn, err := s.config.ClusterArn()
	if err != nil {
		return err
	}
	// ecs run task
	networkConfig, err := s.getNetworkConfig(taskDefArn)
	if err != nil {
		return err
	}

	taskRunOutput, err := s.ecsClient.RunTask(&ecs.RunTaskInput{
		Cluster:              &clusterArn,
		LaunchType:           &launchType,
		TaskDefinition:       &taskDefArn,
		NetworkConfiguration: networkConfig,
	})
	if err != nil {
		return fmt.Errorf("task run failed: %s", err)
	}

	if len(taskRunOutput.Tasks) == 0 {
		return fmt.Errorf("task run failed: Task not found: %s", taskDefArn)
	}

	fmt.Printf("Task %s started\n", taskDefArn)
	if !wait {
		return nil
	}

	var runOutputTaskArns []*string
	for _, taskArn := range taskRunOutput.Tasks {
		runOutputTaskArns = append(runOutputTaskArns, taskArn.TaskArn)
	}

	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterArn,
		Tasks:   runOutputTaskArns,
	}

	// wait for the task to run
	err = s.waitForTask(describeTasksInput)
	if err != nil {
		return err
	}

	// print out the logs
	logEvents, err := s.getLogEvents(taskDefArn, launchType, describeTasksInput)
	if err != nil {
		fmt.Printf("Failed to get logs for %s: %s", taskDefArn, err)
	}
	// TODO compact the print format of log event (currently prints single
	// field per line) to single line per event
	for _, log := range logEvents {
		fmt.Println(log.String())
	}

	return nil
}

func (s *AwsEcs) getNetworkConfig(taskDefArn string) (*ecs.NetworkConfiguration, error) {

	privateSubnets, err := s.config.PrivateSubnets()
	if err != nil {
		return nil, err
	}
	privateSubnetsPt := []*string{}
	for _, subnet := range privateSubnets {
		privateSubnetsPt = append(privateSubnetsPt, &subnet)
	}
	securityGroups, err := s.config.SecurityGroups()
	if err != nil {
		return nil, err
	}
	securityGroupsPt := []*string{}
	for _, subnet := range securityGroups {
		securityGroupsPt = append(securityGroupsPt, &subnet)
	}

	assignPublicIp := "DISABLED"
	awsvpcConfiguration := &ecs.AwsVpcConfiguration{
		AssignPublicIp: &assignPublicIp,
		SecurityGroups: securityGroupsPt,
		Subnets:        privateSubnetsPt,
	}
	networkConfig := &ecs.NetworkConfiguration{
		AwsvpcConfiguration: awsvpcConfiguration,
	}
	return networkConfig, err
}

func (s *AwsEcs) waitForTask(describeTasksInput *ecs.DescribeTasksInput) error {
	err := s.ecsClient.WaitUntilTasksRunning(describeTasksInput)
	if err != nil {
		fmt.Println("Task failed!")
		return err
	}
	result, err := s.ecsClient.DescribeTasks(describeTasksInput)
	taskArn := result.Tasks[0].TaskArn
	if err != nil {
		return err
	}
	container := result.Tasks[0].Containers[0]
	status := container.LastStatus
	if *status != "RUNNING" {
		reason := container.Reason
		fmt.Printf("Container did not start. Current status %s: %s", *status, *reason)
	} else {
		fmt.Printf("Task %s running\n", *taskArn)
		// wait for the task to exit
		err := s.ecsClient.WaitUntilTasksStopped(describeTasksInput)
		if err != nil {
			return err
		}
		fmt.Printf("Task %s stopped\n", *taskArn)
	}

	return nil
}

func (s *AwsEcs) getLogEvents(taskDefArn string, launchType string, describeTasksInput *ecs.DescribeTasksInput) ([]*cloudwatchlogs.OutputLogEvent, error) {

	// get log stream
	result, err := s.ecsClient.DescribeTasks(describeTasksInput)
	if err != nil {
		return nil, err
	}
	container := result.Tasks[0].Containers[0]
	if container.Reason != nil {
		status := container.LastStatus
		reason := container.Reason
		fmt.Printf("Container exited with status %s: %s", *status, *reason)
	}
	logStream := container.RuntimeId

	// get log group
	resultDef, err := s.ecsClient.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefArn,
	})
	if err != nil {
		return nil, err
	}

	taskDef := resultDef.TaskDefinition
	containerDef := taskDef.ContainerDefinitions[0]
	logGroup, ok := containerDef.LogConfiguration.Options["awslogs-group"]
	if !ok {
		fmt.Println("Failed to get logs")
	}
	fmt.Printf("Getting logs for %s, log stream: %s, log group: %s\n", *taskDef.TaskDefinitionArn, *logStream, *logGroup)

	if launchType == "FARGATE" {
		logPrefix := containerDef.LogConfiguration.Options["awslogs-stream-prefix"]
		taskArnSlice := strings.Split(*container.TaskArn, "/")
		taskID := taskArnSlice[len(taskArnSlice)-1]
		stream := fmt.Sprintf("%s/%s/%s", *logPrefix, *containerDef.Name, taskID)
		logStream = &stream
	}

	// get log events
	logsOutput, err := s.logsSrc.awsLogClient.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  logGroup,
		LogStreamName: logStream,
	})
	if err != nil {
		return nil, err
	}

	return logsOutput.Events, nil
}
