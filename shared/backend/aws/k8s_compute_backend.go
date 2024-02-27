package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	dockerterm "github.com/moby/term"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/term"
)

type K8SComputeBackend struct {
	Backend    *Backend
	ClientSet  kubernetes.Interface
	rawConfig  *rest.Config
	KubeConfig kube.K8SConfig
}

const (
	Warning = "Warning"
)

func NewK8SComputeBackend(ctx context.Context, k8sConfig kube.K8SConfig, b *Backend) (*K8SComputeBackend, error) {
	clientset, rawConfig, err := createClientSet(ctx, k8sConfig, b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to instantiate k8s client")
	}

	return &K8SComputeBackend{
		Backend:    b,
		ClientSet:  clientset,
		KubeConfig: k8sConfig,
		rawConfig:  rawConfig,
	}, nil
}

func (k8s *K8SComputeBackend) refreshCredentials() error {
	clientset, rawConfig, err := createClientSet(context.Background(), k8s.KubeConfig, k8s.Backend)
	if err != nil {
		return errors.Wrap(err, "unable to refresh k8s credentials")
	}
	k8s.ClientSet = clientset
	k8s.rawConfig = rawConfig
	return nil
}

func createClientSet(ctx context.Context, k8sConfig kube.K8SConfig, b *Backend) (kubernetes.Interface, *rest.Config, error) {
	return kube.CreateK8sClient(ctx, k8sConfig, kube.AwsClients{
		EksClient:        b.eksclient,
		StsPresignClient: b.stspresignclient,
	}, b.k8sClientCreator)
}

func (k8s *K8SComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	secret, err := k8s.ClientSet.CoreV1().Secrets(k8s.KubeConfig.Namespace).Get(ctx, "integration-secret", v1.GetOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to retrieve integration secret")
	}

	if value, ok := secret.Data["integration_secret"]; ok {
		secret := &config.IntegrationSecret{}
		err = json.Unmarshal(value, secret)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
		}
		arn := ""
		return secret, &arn, nil
	}
	return nil, nil, errors.New("integration_secret key is missing from the integration secret")
}

func (k8s *K8SComputeBackend) GetParam(ctx context.Context, name string) (string, error) {
	configMap, err := k8s.ClientSet.CoreV1().ConfigMaps(k8s.KubeConfig.Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}

	if value, ok := configMap.Data["stacklist"]; ok {
		return value, nil
	}

	return "", errors.Wrapf(err, "unable to retrieve a stacklist key from stacklist configmap")
}

func (k8s *K8SComputeBackend) WriteParam(
	ctx context.Context,
	name string,
	val string,
) error {
	configMap, err := k8s.ClientSet.CoreV1().ConfigMaps(k8s.KubeConfig.Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}
	configMap.Data["stacklist"] = val
	_, err = k8s.ClientSet.CoreV1().ConfigMaps(k8s.KubeConfig.Namespace).Update(ctx, configMap, v1.UpdateOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to update stacklist configmap")
	}
	return nil
}

func (k8s *K8SComputeBackend) GetDeploymentName(stackName string, serviceName string) string {
	return fmt.Sprintf("%s-%s", stackName, serviceName)
}

func (k8s *K8SComputeBackend) PrintLogs(ctx context.Context, stackName, serviceName, containerName string, opts ...util.PrintOption) error {
	if err := k8s.refreshCredentials(); err != nil {
		return errors.Wrap(err, "unable to refresh k8s credentials")
	}

	logrus.Info("***************************************************************")
	logrus.Infof("* Printing logs for stack '%s', service '%s'", stackName, serviceName)
	logrus.Info("***************************************************************")
	deploymentName := k8s.GetDeploymentName(stackName, serviceName)

	pods, err := k8s.getPods(ctx, deploymentName)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve a list of pods")
	}

	logrus.Debugf("Found %d matching pods.", len(pods.Items))
	if len(pods.Items) == 0 {
		return nil
	}

	for _, pod := range pods.Items {
		logrus.Infof("Pod: %s, status: %s", pod.Name, string(pod.Status.Phase))
		logrus.Info("***************************************************************")
		k8s.printPodLogs(ctx, true, pod, pod.Spec.InitContainers, pod.Status.InitContainerStatuses, containerName, opts...)
		k8s.printPodLogs(ctx, false, pod, pod.Spec.Containers, pod.Status.ContainerStatuses, containerName, opts...)
	}

	if k8s.KubeConfig.AuthMethod != kube.AuthMethodEKS {
		return nil
	}

	expression := fmt.Sprintf(`fields @timestamp, log
| sort @timestamp desc
| limit 20
| filter kubernetes.namespace_name = "%s"
| filter kubernetes.pod_name like "%s-%s"
| filter kubernetes.container_name = "%s"`, k8s.KubeConfig.Namespace, stackName, serviceName, containerName)

	logGroup := fmt.Sprintf("/%s/fluentbit-cloudwatch", k8s.KubeConfig.ClusterID)

	logReference := util.LogReference{
		LinkOptions: util.LinkOptions{
			Region:       k8s.Backend.GetAWSRegion(),
			LaunchType:   util.LaunchTypeK8S,
			AWSAccountID: k8s.Backend.GetAWSAccountID(),
		},
		Expression:   expression,
		LogGroupName: logGroup,
	}

	return k8s.Backend.DisplayCloudWatchInsightsLink(ctx, logReference)
}

func (k8s *K8SComputeBackend) promptForContainerName(pod corev1.Pod) (string, error) {
	var containerName string
	containerNames := []string{}
	for _, container := range pod.Spec.Containers {
		containerNames = append(containerNames, container.Name)
	}

	logrus.Warnf("There's more than one container in a pod '%s': %s, and --container flag is not provided.", pod.Name, strings.Join(containerNames, ", "))

	prompt := &survey.Select{
		Message: "Which container are you trying to shell into?",
		Options: containerNames,
		Default: containerNames[0],
	}

	err := survey.AskOne(prompt, &containerName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to ask for a container name")
	}
	if len(containerName) == 0 {
		return "", errors.New("Please specify container name via --container flag")
	}
	return containerName, nil
}

func (k8s *K8SComputeBackend) streamPodLogs(ctx context.Context, pod corev1.Pod, containerName string, follow bool, opts ...util.PrintOption) error {
	logrus.Infof("... streaming logs from pod %s ...", pod.Name)

	opts = append(opts,
		util.WithPaginator(util.NewPodLogPaginator(pod.Name, k8s.ClientSet.CoreV1().Pods(k8s.KubeConfig.Namespace), corev1.PodLogOptions{
			Container: containerName,
			Follow:    follow,
		})),
		util.WithLogTemplate(util.RawStreamMessageTemplate))
	p := util.MakeComputeLogPrinter(ctx, opts...)
	return p.Print(ctx)
}

func (k8s *K8SComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error {
	// Get the cronjob and create a job out of it
	cronJob, err := k8s.ClientSet.BatchV1().CronJobs(k8s.KubeConfig.Namespace).Get(ctx, taskDefArn, v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to retrieve a template cronjob")
	}

	jobId := fmt.Sprintf("%s-%s-job", taskDefArn, uuid.NewUUID())
	jobDef := k8s.createJobFromCronJob(cronJob, jobId)
	jb, err := k8s.ClientSet.BatchV1().Jobs(k8s.KubeConfig.Namespace).Create(ctx, jobDef, v1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to create a k8s job")
	}
	logrus.Infof("Created a k8s job '%s'", jobId)

	timeout := int64(600)
	watch, err := k8s.ClientSet.BatchV1().Jobs(k8s.KubeConfig.Namespace).Watch(ctx, v1.ListOptions{
		FieldSelector:  "metadata.name=" + jobId,
		TimeoutSeconds: &timeout,
	})
	if err != nil {
		return errors.Wrap(err, "unable to start a k8s job watch")
	}

	_, err = k8s.ClientSet.BatchV1().Jobs(k8s.KubeConfig.Namespace).Get(ctx, jobId, v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "k8s job was not created")
	}

	events := watch.ResultChan()
	for {
		event := <-events
		if event.Object == nil {
			return fmt.Errorf("job '%s' result channel closed", jobId)
		}
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return fmt.Errorf("unexpected object type: %T", event.Object)
		}

		terminate := false
		conditions := job.Status.Conditions
		for _, condition := range conditions {
			switch condition.Type {
			case batchv1.JobComplete:
				logrus.Infof("k8s job '%s' completed successfully", jobId)
				terminate = true
			case batchv1.JobFailed:
				return errors.Wrapf(err, "k8s job '%s' failed", job.ObjectMeta.Name)
			default:
				logrus.Debugf("unexpected k8s job '%s' condition: %s", jobId, condition.Type)
				terminate = true
			}
		}
		if terminate {
			break
		}
	}

	logrus.Debug("Interrogating pod jobs")

	podsRef, err := util.IntervalWithTimeout(func() (*corev1.PodList, error) {
		pods, err := k8s.getJobPods(ctx, jb)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to retrieve job pods, will retry")
		}

		if len(pods.Items) == 0 {
			return nil, errors.Wrapf(err, "unable to find any pods created by the job, will retry")
		}

		allPodsReady := true
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodPending {
				allPodsReady = false
				break
			}
		}
		if !allPodsReady {
			return nil, errors.Wrapf(err, "not all pods are ready")
		}

		return pods, nil
	}, 10*time.Second, 2*time.Minute)

	if err != nil {
		return errors.Wrap(err, "pods did not successfuly start on time")
	}

	if podsRef == nil {
		return errors.New("nil reference")
	}

	pods := *podsRef

	if pods == nil || len(pods.Items) == 0 {
		return errors.New("No pods found for this migration job, the job most likely failed.")
	}

	logrus.Debugf("Found %d successfuly started pods", len(pods.Items))

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			logrus.Debugf("Pod: %s, container %s, status: %s", pod.Name, container.Name, pod.Status.Phase)
			err = k8s.streamPodLogs(ctx, pod, container.Name, true)
			if err != nil {
				logrus.Error(err.Error())
			}
		}
	}

	// Delete the job
	policy := v1.DeletePropagationBackground
	err = k8s.ClientSet.BatchV1().Jobs(k8s.KubeConfig.Namespace).Delete(ctx, jb.Name, v1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	return err
}

func (k8s *K8SComputeBackend) getPods(ctx context.Context, deploymentName string) (*corev1.PodList, error) {
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	return k8s.getSelectorPods(ctx, labelSelector)
}

func (k8s *K8SComputeBackend) getTargetGroupBindings(ctx context.Context, stackName string, serviceName string) (*unstructured.UnstructuredList, error) {
	dynamic := dynamic.NewForConfigOrDie(k8s.rawConfig)

	gvk := schema.FromAPIVersionAndKind("elbv2.k8s.aws/v1beta1", "TargetGroupBinding")
	gv := gvk.GroupVersion()
	target := gv.WithResource("targetgroupbindings")

	deploymentName := k8s.GetDeploymentName(stackName, serviceName)
	labelSelector := fields.SelectorFromSet(fields.Set{
		"ingress.k8s.aws/stack": fmt.Sprintf("service-%s", deploymentName),
	})

	// Side note: CRD field selectors are not supported
	resources, err := dynamic.Resource(target).Namespace(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{Limit: 100000,
		LabelSelector: labelSelector.String()})

	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve a list of target group bindings for deployment")
	}

	return resources, nil
}

func (k8s *K8SComputeBackend) getJobPods(ctx context.Context, job *batchv1.Job) (*corev1.PodList, error) {
	labelSelector, err := v1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse a secelector: %v", job.Spec.Selector)
	}
	pods, err := k8s.ClientSet.CoreV1().Pods(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector.String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve a list of pods for selector %s", labelSelector.String())
	}
	return pods, nil
}

func (k8s *K8SComputeBackend) getSelectorPods(ctx context.Context, labelSelector v1.LabelSelector) (*corev1.PodList, error) {
	selector := labels.Set(labelSelector.MatchLabels).String()
	pods, err := k8s.ClientSet.CoreV1().Pods(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve a list of pods for selector %s", selector)
	}
	return pods, nil
}

func (k8s *K8SComputeBackend) Shell(ctx context.Context, stackName, serviceName, containerName string, shellCommand []string) error {
	deploymentName := k8s.GetDeploymentName(stackName, serviceName)

	pods, err := k8s.getPods(ctx, deploymentName)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve a list of pods")
	}

	if len(pods.Items) == 0 {
		return errors.New("No matching pods found")
	}

	podName := pods.Items[0].Name

	pod, err := k8s.ClientSet.CoreV1().Pods(k8s.KubeConfig.Namespace).Get(ctx, podName, v1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve pod information for %s", podName)
	}

	if len(pod.Spec.Containers) > 1 && len(containerName) == 0 {
		if diagnostics.IsInteractiveContext(ctx) {
			var err error
			containerName, err = k8s.promptForContainerName(*pod)
			if err != nil {
				return errors.Wrap(err, "failed to prompt for container name")
			}
		}
	}

	if len(containerName) == 0 {
		containerName = pod.Spec.Containers[0].Name
	}

	logrus.Debugf("Found %d matching pods. Opening a TTY tunnel into pod '%s', container '%s', command '%s'", len(pods.Items), podName, containerName, shellCommand)

	req := k8s.ClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(pod.Name).Namespace(pod.Namespace).SubResource("exec").Param("container", containerName)

	cmd := strings.Fields("sh -c /bin/sh")
	if len(shellCommand) > 0 {
		cmd = append(strings.Fields("sh -c"), shellCommand...)
	}
	eo := &corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdout:    true,
		Stdin:     true,
		Stderr:    false,
		TTY:       true,
	}

	// Non-interactive shell
	if !diagnostics.IsInteractiveContext(ctx) {
		eo = &corev1.PodExecOptions{
			Container: containerName,
			Command:   shellCommand,
			Stdin:     false,
			Stdout:    true,
			Stderr:    false,
			TTY:       false,
		}
	}

	req.VersionedParams(eo, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k8s.rawConfig, http.MethodPost, req.URL())
	if err != nil {
		return errors.Wrapf(err, "unable to create a SPDY executor")
	}
	stdin, stdout, stderr := dockerterm.StdStreams()

	streamOptions := remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	}

	t := term.TTY{
		In:  stdin,
		Out: stdout,
		Raw: true,
	}

	if !diagnostics.IsInteractiveContext(ctx) {
		streamOptions = remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdout,
			Tty:    false,
		}
		t = term.TTY{
			Out: stdout,
			Raw: true,
		}
	}

	streamOptions.TerminalSizeQueue = t.MonitorSize(t.GetSize())
	return t.Safe(func() error { return exec.StreamWithContext(ctx, streamOptions) })
}

// This function is used to retrieve events for a given stack, it looks into Deployment, Pod, Ingress, HorizontalPodAutoscaler and TargetGroupBinding triggered events
func (k8s *K8SComputeBackend) GetEvents(ctx context.Context, stackName string, services []string) error {
	if len(services) == 0 {
		return errors.Errorf("No services are defined for stack '%s'", stackName)
	}

	eventsFound := false

	for _, serviceName := range services {
		resourceEvents, err := k8s.getServiceEvents(ctx, stackName, serviceName)
		if err != nil {
			return errors.Wrapf(err, "unable to retrieve events for service '%s'", serviceName)
		}

		k8s.interpretEvents(stackName, serviceName, resourceEvents)
		eventsFound = eventsFound || len(resourceEvents) > 0
	}

	if !eventsFound {
		logrus.Info("No events found for this stack")
	}

	return nil
}

func (k8s *K8SComputeBackend) getServiceEvents(ctx context.Context, stackName string, serviceName string) ([]corev1.Event, error) {
	resourceEvents := make([]corev1.Event, 0)
	deploymentName := k8s.GetDeploymentName(stackName, serviceName)

	// Get all pods in a deployment
	pods, err := k8s.getPods(ctx, deploymentName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve a list of pods")
	}
	if len(pods.Items) == 0 {
		return nil, errors.New("No matching pods found, unable to retrieve events")
	}

	// Get events for every pods in a deployment
	for _, pod := range pods.Items {
		events, err := k8s.getResourceEvents(ctx, pod.Name, "Pod")
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve events for a pod")
		}
		resourceEvents = append(resourceEvents, events.Items...)
		for _, status := range pod.Status.ContainerStatuses {
			if status.RestartCount > 0 {
				resourceEvents = append(resourceEvents, corev1.Event{
					InvolvedObject: corev1.ObjectReference{
						Name: pod.Name,
						Kind: "Pod",
					},
					Reason:  "HappyRestartCount",
					Type:    Warning,
					Message: fmt.Sprintf("Container %s in pod %s restarted %d times", status.Name, pod.Name, status.RestartCount),
				})
			}

			if status.LastTerminationState.Terminated != nil && status.LastTerminationState.Terminated.ExitCode != 0 {
				resourceEvents = append(resourceEvents, corev1.Event{
					InvolvedObject: corev1.ObjectReference{
						Name: pod.Name,
						Kind: "Pod",
					},
					Reason:  "HappyTerminated",
					Type:    Warning,
					Message: fmt.Sprintf("Container %s in pod %s exited with code %d", status.Name, pod.Name, status.LastTerminationState.Terminated.ExitCode),
				})
			}
		}
	}

	// Get events for the deployment, skipping ReplicaSet events for now
	events, err := k8s.getResourceEvents(ctx, deploymentName, "Deployment")
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve events for a deployment")
	}
	resourceEvents = append(resourceEvents, events.Items...)

	// Get events for the service
	events, err = k8s.getResourceEvents(ctx, deploymentName, "Service")
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve events for a service")
	}
	resourceEvents = append(resourceEvents, events.Items...)

	// Get events for the horizontal pod autoscaler
	events, err = k8s.getResourceEvents(ctx, deploymentName, "HorizontalPodAutoscaler")
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve events for a horizontal pod autoscaler")
	}
	resourceEvents = append(resourceEvents, events.Items...)

	// Get events for the ingress
	events, err = k8s.getResourceEvents(ctx, deploymentName, "Ingress")
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve events for an ingress resource")
	}
	resourceEvents = append(resourceEvents, events.Items...)

	// Find all matching target group bindings, and events for them. Target groups are created from the ingress resource. Target
	// groups are labeled with the "ingress.k8s.aws/stack=service-<STACK_NAME>-<SERVICE_NAME>" label.
	targetGroupBindings, err := k8s.getTargetGroupBindings(ctx, stackName, serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve a list of ALB target group bindings")
	}

	// Get events for all target group bindings
	for _, targetGroupBinding := range targetGroupBindings.Items {
		events, err = k8s.getResourceEvents(ctx, targetGroupBinding.GetName(), "TargetGroupBinding")
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve events for a target group binding")
		}
		resourceEvents = append(resourceEvents, events.Items...)
	}

	// Sort events by timestamp
	if len(resourceEvents) > 1 {
		sort.Slice(resourceEvents, func(i, j int) bool {
			if resourceEvents[i].FirstTimestamp.Equal(&resourceEvents[j].FirstTimestamp) {
				return resourceEvents[i].InvolvedObject.Name < resourceEvents[j].InvolvedObject.Name
			}
			return resourceEvents[i].FirstTimestamp.Before(&resourceEvents[j].FirstTimestamp)
		})
	}
	return resourceEvents, nil
}

func (k8s *K8SComputeBackend) Describe(ctx context.Context, stackName string, serviceName string) (interfaces.StackServiceDescription, error) {
	params := make(map[string]string)
	params["namespace"] = k8s.KubeConfig.Namespace
	params["deployment_name"] = k8s.GetDeploymentName(stackName, serviceName)
	params["auth_method"] = k8s.KubeConfig.AuthMethod
	params["kube_api"] = k8s.rawConfig.Host

	description := interfaces.StackServiceDescription{
		Compute: "K8S",
		Params:  params,
	}
	return description, nil
}

func (k8s *K8SComputeBackend) createJobFromCronJob(cronJob *batchv1.CronJob, jobName string) *batchv1.Job {
	annotations := map[string]string{
		"cronjob.kubernetes.io/instantiate": "manual",
	}
	for k, v := range cronJob.Spec.JobTemplate.Annotations {
		annotations[k] = v
	}

	var ttl int32 = 600
	cronJob.Spec.JobTemplate.Spec.TTLSecondsAfterFinished = &ttl
	job := &batchv1.Job{
		TypeMeta: v1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: v1.ObjectMeta{
			Name:        jobName,
			Annotations: annotations,
			Labels:      cronJob.Spec.JobTemplate.Labels,
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion: batchv1.SchemeGroupVersion.String(),
					Kind:       "CronJob",
					Name:       cronJob.GetName(),
					UID:        cronJob.GetUID(),
				},
			},
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	job.Namespace = k8s.KubeConfig.Namespace

	return job
}

func (k8s *K8SComputeBackend) GetResources(ctx context.Context, stackName string) ([]util.ManagedResource, error) {
	managedResources := []util.ManagedResource{}
	dynamic := dynamic.NewForConfigOrDie(k8s.rawConfig)
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(k8s.rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable create discovery client")
	}
	groupResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, errors.Wrap(err, "unable to discover preferred resource versions")
	}

	// Code below can be used to enumerate all resources in a K8s cluster.
	// Some resources are transient and we don't want to discover them. Here are some examples:
	// * TargetGroupBinding
	// * PodMetrics
	// * EndpointSlice
	// * Event
	// * ReplicaSet
	// * Pod
	// * ControllerRevision
	// etc...

	beKind := map[string]bool{
		"Deployment":              true,
		"Service":                 true,
		"ConfigMap":               true,
		"Secret":                  true,
		"Ingress":                 true,
		"PersistentVolumeClaim":   true,
		"StatefulSet":             true,
		"DaemonSet":               true,
		"Job":                     true,
		"CronJob":                 true,
		"HorizontalPodAutoscaler": true,
		"PodDisruptionBudget":     true,
		"NetworkPolicy":           true,
		"Role":                    true,
		"RoleBinding":             true,
		"ServiceAccount":          true,
	}
	for _, gr := range groupResources {
		for _, resource := range gr.APIResources {
			if _, ok := beKind[resource.Kind]; !ok {
				continue
			}
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, resource.Kind)
			gv := gvk.GroupVersion()
			target := gv.WithResource(resource.Name)

			resourceMap := map[types.UID]bool{}

			// These resources are guaranteed as they are matched by a label selector
			ls := &v1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/part-of": stackName}}

			resources, err := dynamic.Resource(target).Namespace(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
				LabelSelector: v1.FormatLabelSelector(ls),
			})
			if err != nil {
				// This is an expected behavior, if resources do not exist, we just skip them
				logrus.Debugf("unable to retrieve a list of resources %s/%s in namespace %s: %s", resource.Kind, resource.Name, k8s.KubeConfig.Namespace, err.Error())
				continue
			}

			for _, item := range resources.Items {
				resourceMap[item.GetUID()] = true
				managedResources = append(managedResources, util.ManagedResource{
					ManagedBy: "k8s",
					Name:      item.GetName(),
					Type:      item.GetKind(),
					Provider:  "k8s",
					Module:    "",
					Instances: []string{},
				})

				if resource.Kind == "Ingress" {
					managedResources = append(managedResources, extractIngressResources(item)...)
				}
			}

			resources, err = dynamic.Resource(target).Namespace(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{})
			if err != nil {
				logrus.Errorf("unable to retrieve a list of resources %s/%s in namespace %s: %s", resource.Kind, resource.Name, k8s.KubeConfig.Namespace, err.Error())
				continue
			}

			// Resources that have not been identified by app.kubernetes.io/part-of label (e.g. created by an older version of the happy-stack-eks module),
			// these are most likely still managed by a stack, but that's not guaranteed. Added for compatibility.
			for _, item := range resources.Items {
				if _, ok := resourceMap[item.GetUID()]; ok {
					continue // Already accounted for
				}
				resourceMap[item.GetUID()] = true
				if strings.Index(item.GetName(), fmt.Sprintf("%s-", stackName)) != 0 {
					continue
				}

				managedResources = append(managedResources, util.ManagedResource{
					ManagedBy: "k8s",
					Name:      "*" + item.GetName(),
					Type:      item.GetKind(),
					Provider:  "k8s",
					Module:    "",
					Instances: []string{},
				})

				if resource.Kind == "Ingress" {
					managedResources = append(managedResources, extractIngressResources(item)...)
				}
			}
		}
	}

	return managedResources, nil
}

func extractIngressResources(item unstructured.Unstructured) []util.ManagedResource {
	managedResources := []util.ManagedResource{}
	value, found, err := unstructured.NestedSlice(item.Object, "status", "loadBalancer", "ingress")
	if err == nil {
		if found {
			for _, v := range value {
				if lbIngress, ok := v.(map[string]interface{}); ok {
					if lbHost := lbIngress["hostname"]; lbHost != nil {
						managedResources = append(managedResources, util.ManagedResource{
							ManagedBy: "k8s",
							Name:      "",
							Type:      "Application Load Balancer",
							Provider:  "ALB Ingress Controller",
							Module:    fmt.Sprintf("k8s:%s/%s", item.GetKind(), item.GetName()),
							Instances: []string{lbHost.(string)},
						})
					}
				}
			}
		}
	}

	if rules, ok, err := unstructured.NestedSlice(item.Object, "spec", "rules"); ok && err == nil {
		for i := range rules {
			rule, ok := rules[i].(map[string]interface{})
			if !ok {
				continue
			}
			host := rule["host"]
			if host != nil && host != "" {
				managedResources = append(managedResources, util.ManagedResource{
					ManagedBy: "k8s",
					Name:      "",
					Type:      "Route 53 Entry",
					Provider:  "External DNS",
					Module:    fmt.Sprintf("k8s:%s/%s", item.GetKind(), item.GetName()),
					Instances: []string{host.(string)},
				})
			}
		}
	}
	return managedResources
}

func (k8s *K8SComputeBackend) getResourceEvents(ctx context.Context, resourceName string, resourceKind string) (*corev1.EventList, error) {
	fieldSelector := fields.SelectorFromSet(fields.Set{
		"involvedObject.name": resourceName,
		"involvedObject.kind": resourceKind,
	})

	events, err := k8s.ClientSet.CoreV1().Events(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
		FieldSelector: fieldSelector.String(),
		Limit:         100,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve events for resource %s/%s", resourceKind, resourceName)
	}
	return events, nil
}

func (k8s *K8SComputeBackend) interpretEvents(stackName string, serviceName string, events []corev1.Event) {

	messages := []string{}

	for _, e := range events {
		if e.Type == Warning {
			logrus.Warnf("%s/%s - %s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message)

			for _, signal := range K8sDebugSignals {
				if e.Reason == signal.Reason && e.InvolvedObject.Kind == signal.Kind && strings.Contains(e.Message, signal.MessageSignature) {
					var sb strings.Builder

					if logrus.GetLevel() == logrus.DebugLevel {
						sb.WriteString(e.InvolvedObject.Kind)
						sb.WriteString("/")
						sb.WriteString(e.InvolvedObject.Name)
						sb.WriteString(": ")
					}
					sb.WriteString(signal.Description)

					if logrus.GetLevel() == logrus.DebugLevel {
						sb.WriteString(" [")
						sb.WriteString(e.Message)
						sb.WriteString("] ")
						if signal.Remediation != "" {
							sb.WriteString(" -- ")
							sb.WriteString(signal.Remediation)
						}
						sb.WriteString(" -- ")
						sb.WriteString("See ")
						sb.WriteString(signal.RunbookUrl)
					}

					messages = append(messages, sb.String())
				}
			}
		} else {
			logrus.Infof("%s/%s - %s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message)
		}
	}

	if len(messages) > 0 {
		logrus.Println()
		logrus.Println("Many \"Warning\" events - please check to see whether your service is crashing:")
		logrus.Infof("  happy --env %s logs %s %s", k8s.Backend.Conf().GetEnv(), stackName, serviceName)
		logrus.Println()
		logrus.Infof("Here's a list of issues we've detected for service '%s' in stack '%s':", serviceName, stackName)
		deduper := map[string]bool{}
		for _, m := range messages {
			if _, ok := deduper[m]; !ok {
				deduper[m] = true
				logrus.Warn(m)
			}
		}
		logrus.Println()
	}

}

func (k8s *K8SComputeBackend) ListHappyNamespaces(ctx context.Context) ([]string, error) {
	happyNamespaces := []string{}
	namespaces, err := k8s.ClientSet.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list namespaces")
	}
	for _, namespace := range namespaces.Items {
		if namespace.Status.Phase != corev1.NamespaceActive {
			continue
		}
		secrets, err := k8s.ClientSet.CoreV1().Secrets(namespace.Name).List(ctx, v1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to list secrets in namespace %s", namespace.Name)
		}

		for _, secret := range secrets.Items {
			if secret.Type == corev1.SecretTypeOpaque && secret.Name == "integration-secret" {
				happyNamespaces = append(happyNamespaces, namespace.Name)
				break
			}
		}
	}
	return happyNamespaces, nil
}

func (k8s *K8SComputeBackend) GetSecret(ctx context.Context, name string) (map[string][]byte, error) {
	secret, err := k8s.ClientSet.CoreV1().Secrets(k8s.KubeConfig.Namespace).
		Get(ctx, name, v1.GetOptions{})
	if err != nil && !k8serr.IsNotFound(err) {
		return nil, errors.Wrapf(err, "unable to retrieve secret [%s]", name)
	}

	return secret.Data, nil
}

func (k8s *K8SComputeBackend) WriteKeyToSecret(ctx context.Context, name, key, val string, labels map[string]string) (map[string][]byte, error) {
	// make sure the secret exists
	err := k8s.CreateSecretIfNotExists(ctx, name, labels)
	if err != nil {
		return nil, err
	}

	// set the key in secret to the specified value
	patchSecret := corev1.Secret{
		Data: map[string][]byte{
			key: []byte(val),
		},
	}
	patchData, err := json.Marshal(patchSecret)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to marshal secret [%s]", key)
	}
	_, err = k8s.ClientSet.CoreV1().Secrets(k8s.KubeConfig.Namespace).
		Patch(ctx, name, types.StrategicMergePatchType, patchData, v1.PatchOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to set secret [%s]", key)
	}

	return patchSecret.Data, nil
}

func (k8s *K8SComputeBackend) DeleteKeyFromSecret(ctx context.Context, name, key string, labels map[string]string) error {
	// make sure the secret exists
	err := k8s.CreateSecretIfNotExists(ctx, name, labels)
	if err != nil {
		return err
	}

	_, err = k8s.ClientSet.CoreV1().Secrets(k8s.KubeConfig.Namespace).
		Patch(ctx, name, types.JSONPatchType, []byte(fmt.Sprintf("[{\"op\": \"remove\", \"path\": \"/data/%s\"}]", key)), v1.PatchOptions{})
	// if the specified path doesn't exist in the secret then an "invalid" error is returned, which we can ignore
	if err != nil && !k8serr.IsInvalid(err) {
		return errors.Wrapf(err, "unable to set secret [%s]", key)
	}

	return nil
}

func (k8s *K8SComputeBackend) CreateSecretIfNotExists(ctx context.Context, name string, labels map[string]string) error {
	_, err := k8s.ClientSet.CoreV1().Secrets(k8s.KubeConfig.Namespace).
		Create(ctx, &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:   name,
				Labels: labels,
			},
		}, v1.CreateOptions{})
	if err != nil && !k8serr.IsAlreadyExists(err) {
		return errors.Wrapf(err, "unable to create secret [%s]", name)
	}
	return nil
}

func (k8s *K8SComputeBackend) containerStateToString(status *corev1.ContainerStatus) string {
	if status == nil {
		return "unknown"
	}
	if status.State.Running != nil {
		return "running"
	}
	if status.State.Waiting != nil {
		return fmt.Sprintf("waiting (Reason: %s, Message: %s)", status.State.Waiting.Reason, status.State.Waiting.Message)
	}
	if status.State.Terminated != nil {
		return fmt.Sprintf("terminated (Code: %d, Reason: %s, Message: %s)", status.State.Terminated.ExitCode, status.State.Terminated.Reason, status.State.Terminated.Message)
	}

	return fmt.Sprintf("%v", status)
}

func (k8s *K8SComputeBackend) findContainerStatus(statuses []corev1.ContainerStatus, containerName string) *corev1.ContainerStatus {
	for _, status := range statuses {
		if status.Name == containerName {
			return &status
		}
	}
	return nil
}

func (k8s *K8SComputeBackend) printPodLogs(ctx context.Context, init bool, pod corev1.Pod, containers []corev1.Container, containerStatuses []corev1.ContainerStatus, containerName string, opts ...util.PrintOption) {
	for _, container := range containers {
		if containerName == "" || container.Name == containerName {
			containerStatus := k8s.findContainerStatus(containerStatuses, container.Name)
			if containerStatus == nil {
				continue
			}

			restartCount := 0
			status := "unknown"
			healthy := false

			restartCount = int(containerStatus.RestartCount)
			status = k8s.containerStateToString(containerStatus)
			if init {
				if containerStatus.State.Running != nil {
					// If the init container is running, and never restarted, it's healthy
					healthy = restartCount == 0
				} else if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode == 0 {
					// If the init container is terminated, and it exited successfully, it's healthy
					healthy = true
				}
			} else {
				// For regular containers, it's healthy if it's running and never restarted
				healthy = restartCount == 0 && containerStatus.State.Running != nil
			}

			label := "[container]"
			if init {
				label = "[init_container]"
			}

			logrus.Info("---------------------------------------------------------------------")
			if !healthy {
				logrus.Errorf("%s '%s' in pod '%s' is not healthy. It's in '%s' status, and it restarted %d times.", label, container.Name, pod.Name, status, restartCount)
			}
			logrus.Info("---------------------------------------------------------------------")

			err := k8s.streamPodLogs(ctx, pod, container.Name, false, opts...)
			if err != nil {
				logrus.Error(err.Error())
			}
		}
	}
}
