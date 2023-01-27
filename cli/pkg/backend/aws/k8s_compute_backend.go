package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/cli/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/util"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	dockerterm "github.com/moby/term"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/uuid"
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

func NewK8SComputeBackend(ctx context.Context, k8sConfig kube.K8SConfig, b *Backend, clientCreator kube.K8sClientCreator) (interfaces.ComputeBackend, error) {
	clientset, rawConfig, err := kube.CreateK8sClient(ctx, k8sConfig, kube.AwsClients{
		EksClient: b.eksclient,
		StsClient: b.stspresignclient,
	}, clientCreator)

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

func (k8s *K8SComputeBackend) getDeploymentName(stackName string, serviceName string) string {
	return fmt.Sprintf("%s-%s", stackName, serviceName)
}

func (k8s *K8SComputeBackend) PrintLogs(ctx context.Context, stackName string, serviceName string, opts ...util.PrintOption) error {
	pods, err := k8s.getPods(ctx, stackName, serviceName)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve a list of pods")
	}

	logrus.Infof("Found %d matching pods.", len(pods.Items))

	for _, pod := range pods.Items {
		err = k8s.streamPodLogs(ctx, pod, false, opts...)
		if err != nil {
			logrus.Error(err.Error())
		}
	}
	return nil
}

func (k8s *K8SComputeBackend) streamPodLogs(ctx context.Context, pod corev1.Pod, follow bool, opts ...util.PrintOption) error {
	logrus.Infof("... streaming logs from pod %s ...", pod.Name)

	opts = append(opts,
		util.WithPaginator(util.NewPodLogPaginator(pod.Name, k8s.ClientSet.CoreV1().Pods(k8s.KubeConfig.Namespace), corev1.PodLogOptions{
			Follow: follow,
		})),
		util.WithLogTemplate(util.RawStreamMessageTemplate))
	p := util.MakeComputeLogPrinter(ctx, opts...)
	return p.Print(ctx)
}

func (k8s *K8SComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType config.LaunchType) error {
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

	logrus.Debug("Waiting for all the pods to start")

	podsRef, err := util.IntervalWithTimeout(func() (*corev1.PodList, error) {
		pods, err := k8s.getJobPods(ctx, jb.Name)
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
	}, 10*time.Second, 5*time.Minute)

	if err != nil {
		return errors.Wrap(err, "pods did not successfuly start on time")
	}

	if podsRef == nil {
		return errors.New("nil reference")
	}

	pods := *podsRef

	if pods == nil {
		return errors.New("nil pods reference")
	}

	logrus.Debugf("Found %d successfuly started pods", len(pods.Items))
	for _, pod := range pods.Items {
		err = k8s.streamPodLogs(ctx, pod, true)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	// Delete the job
	policy := v1.DeletePropagationBackground
	err = k8s.ClientSet.BatchV1().Jobs(k8s.KubeConfig.Namespace).Delete(ctx, jb.Name, v1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	return err
}

func (k8s *K8SComputeBackend) getPods(ctx context.Context, stackName string, serviceName string) (*corev1.PodList, error) {
	deploymentName := k8s.getDeploymentName(stackName, serviceName)
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	return k8s.getSelectorPods(ctx, labelSelector)
}

func (k8s *K8SComputeBackend) getJobPods(ctx context.Context, taskDefArn string) (*corev1.PodList, error) {
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"job-name": taskDefArn}}
	return k8s.getSelectorPods(ctx, labelSelector)
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

func (k8s *K8SComputeBackend) Shell(ctx context.Context, stackName string, serviceName string) error {
	pods, err := k8s.getPods(ctx, stackName, serviceName)
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

	if len(pod.Spec.Containers) > 1 {
		return errors.Errorf("There's more than one container in a pod '%s'", podName)
	}

	containerName := pod.Spec.Containers[0].Name

	logrus.Infof("Found %d matching pods. Opening a TTY tunnel into pod '%s', container '%s'", len(pods.Items), podName, containerName)

	req := k8s.ClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(pod.Name).Namespace(pod.Namespace).SubResource("exec").Param("container", containerName)

	eo := &corev1.PodExecOptions{
		Container: containerName,
		Command:   strings.Fields("sh -c /bin/sh"),
		Stdout:    true,
		Stdin:     true,
		Stderr:    false,
		TTY:       true,
	}
	req.VersionedParams(eo, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k8s.rawConfig, http.MethodPost, req.URL())
	if err != nil {
		return errors.Wrapf(err, "Unable to create a SPDY executor")
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
	streamOptions.TerminalSizeQueue = t.MonitorSize(t.GetSize())
	return t.Safe(func() error { return exec.StreamWithContext(ctx, streamOptions) })
}

func (k8s *K8SComputeBackend) GetEvents(ctx context.Context, stackName string, services []string) error {
	if len(services) == 0 {
		return nil
	}

	for _, serviceName := range services {
		pods, err := k8s.getPods(ctx, stackName, serviceName)
		if err != nil {
			return errors.Wrap(err, "unable to retrieve a list of pods")
		}

		warnings := make([]corev1.Event, 0)

		deploymentName := k8s.getDeploymentName(stackName, serviceName)
		fieldSelector := fields.SelectorFromSet(fields.Set{
			"involvedObject.name": deploymentName,
			"type":                Warning,
		})

		events, _ := k8s.ClientSet.CoreV1().Events(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
			FieldSelector: fieldSelector.String(),
			TypeMeta:      v1.TypeMeta{Kind: "Deployment"},
		})

		warnings = append(warnings, events.Items...)

		for _, pod := range pods.Items {
			fieldSelector := fields.SelectorFromSet(fields.Set{
				"involvedObject.name": pod.Name,
				"type":                Warning,
			})

			events, _ := k8s.ClientSet.CoreV1().Events(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
				FieldSelector: fieldSelector.String(),
				TypeMeta:      v1.TypeMeta{Kind: "Pod"},
			})

			warnings = append(warnings, events.Items...)
		}

		if len(warnings) > 1 {
			sort.Slice(warnings, func(i, j int) bool {
				if warnings[i].FirstTimestamp.Equal(&warnings[j].FirstTimestamp) {
					return warnings[i].InvolvedObject.Name < warnings[j].InvolvedObject.Name
				}
				return warnings[i].FirstTimestamp.Before(&warnings[j].FirstTimestamp)
			})
		}

		for _, e := range warnings {
			logrus.Warnf("%s/%s - %s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message)
			warnings = append(warnings, e)
		}

		if len(events.Items) >= 1 {
			logrus.Println()
			logrus.Println("Many \"Warning\" events - please check to see whether your service is crashing:")
			logrus.Infof("  happy --env %s logs %s %s", k8s.Backend.Conf().GetEnv(), stackName, serviceName)
		}
	}

	return nil
}

func (k8s *K8SComputeBackend) Describe(ctx context.Context, stackName string, serviceName string) (interfaces.StackServiceDescription, error) {
	params := make(map[string]string)
	params["namespace"] = k8s.KubeConfig.Namespace
	params["deployment_name"] = k8s.getDeploymentName(stackName, serviceName)
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
