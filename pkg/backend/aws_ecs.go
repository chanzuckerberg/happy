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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TaskType string

const (
	DeletionTask  TaskType = "delete"
	MigrationTask TaskType = "migrate"
)

type TaskRunner interface {
	RunTask(taskDef string, launchType string) error
	GetECSClient() ecsiface.ECSAPI
	GetEC2Client() *ec2.EC2
}

type AwsEcs struct {
	session   *session.Session
	ecsClient ecsiface.ECSAPI
	awsConfig *aws.Config
	config    config.HappyConfig
	logsSrc   *Cloudwatchlogs
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

func GetAwsEcs(config config.HappyConfig) TaskRunner {
	awsProfile := config.AwsProfile()
	creatECSMOnce.Do(func() {
		awsConfig := &aws.Config{
			// TODO(el): don't hard-code region
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

func (s *AwsEcs) RunTask(taskDefArn string, launchType string) error {
	log.Printf("Running tasks for %s\n", taskDefArn)

	clusterArn := s.config.ClusterArn()

	// ecs run task
	networkConfig := s.getNetworkConfig(taskDefArn)

	taskRunOutput, err := s.ecsClient.RunTask(&ecs.RunTaskInput{
		Cluster:              &clusterArn,
		LaunchType:           &launchType,
		TaskDefinition:       &taskDefArn,
		NetworkConfiguration: networkConfig,
	})
	if err != nil {
		return errors.Errorf("task run failed: %s", err)
	}

	if len(taskRunOutput.Tasks) == 0 {
		return errors.Errorf("task run failed: Task not found: %s", taskDefArn)
	}

	log.Printf("Task %s started\n", taskDefArn)

	var runOutputTaskArns []*string
	for _, taskArn := range taskRunOutput.Tasks {
		runOutputTaskArns = append(runOutputTaskArns, taskArn.TaskArn)
	}

	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterArn,
		Tasks:   runOutputTaskArns,
	}

	// for the task to run
	err = s.waitForTask(describeTasksInput)
	if err != nil {
		return err
	}

	// print out the logs
	logEvents, err := s.getLogEvents(taskDefArn, launchType, describeTasksInput)
	if err != nil {
		log.Errorf("Failed to get logs for %s: %s", taskDefArn, err)
	}
	// TODO compact the print format of log event (currently prints single
	// field per line) to single line per event
	log.Println("\n\nLOG:")
	log.Println("************************************")
	for _, logEvent := range logEvents {
		log.Println(*logEvent.Message)
	}
	log.Println("************************************")
	return nil
}

func (s *AwsEcs) getNetworkConfig(taskDefArn string) *ecs.NetworkConfiguration {
	privateSubnets := s.config.PrivateSubnets()
	privateSubnetsPt := []*string{}
	for _, subnet := range privateSubnets {
		privateSubnetsPt = append(privateSubnetsPt, &subnet)
	}
	securityGroups := s.config.SecurityGroups()
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
	return networkConfig
}

func (s *AwsEcs) waitForTask(describeTasksInput *ecs.DescribeTasksInput) error {
	err := s.ecsClient.WaitUntilTasksRunning(describeTasksInput)
	if err != nil {
		return errors.Wrap(err, "task failed")
	}
	result, err := s.ecsClient.DescribeTasks(describeTasksInput)
	taskArn := result.Tasks[0].TaskArn
	if err != nil {
		return errors.Wrap(err, "could not describe task")
	}
	container := result.Tasks[0].Containers[0]
	status := container.LastStatus
	if *status != "RUNNING" {
		reason := container.Reason
		log.Warnf("Container did not start. Current status %s: %s", *status, *reason)
	} else {
		log.Printf("Task %s running\n", *taskArn)
		// wait for the task to exit
		err := s.ecsClient.WaitUntilTasksStopped(describeTasksInput)
		if err != nil {
			return err
		}
		log.Printf("Task %s stopped\n", *taskArn)
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
		log.Warnf("Container exited with status %s: %s", *status, *reason)
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
		log.Warnln("Failed to get logs")
	}
	log.Printf("Getting logs for %s, log stream: %s, log group: %s\n", *taskDef.TaskDefinitionArn, *logStream, *logGroup)

	if launchType == config.LaunchTypeFargate {
		logPrefix, ok := containerDef.LogConfiguration.Options["awslogs-stream-prefix"]
		if !ok || logPrefix == nil {
			return nil, errors.New("failed to get a log prefix")
		}
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
