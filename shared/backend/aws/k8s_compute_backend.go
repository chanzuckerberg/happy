package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	dockerterm "github.com/moby/term"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

func NewK8SComputeBackend(ctx context.Context, k8sConfig kube.K8SConfig, b *Backend) (interfaces.ComputeBackend, error) {
	clientset, rawConfig, err := kube.CreateK8sClient(ctx, k8sConfig, kube.AwsClients{
		EksClient:        b.eksclient,
		StsPresignClient: b.stspresignclient,
	}, b.k8sClientCreator)

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
	deploymentName := k8s.getDeploymentName(stackName, serviceName)

	pods, err := k8s.getPods(ctx, deploymentName)
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

	if k8s.KubeConfig.AuthMethod != kube.AuthMethodEKS {
		return nil
	}

	expression := fmt.Sprintf(`fields @timestamp, log
| sort @timestamp desc
| limit 20
| filter kubernetes.namespace_name = "%s"
| filter kubernetes.pod_name like "%s-%s"`, k8s.KubeConfig.Namespace, stackName, serviceName)

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

func (k8s *K8SComputeBackend) getPods(ctx context.Context, deploymentName string) (*corev1.PodList, error) {
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	return k8s.getSelectorPods(ctx, labelSelector)
}

func (k8s *K8SComputeBackend) getTargetGroupBindings(ctx context.Context, stackName string, serviceName string) (*unstructured.UnstructuredList, error) {
	dynamic := dynamic.NewForConfigOrDie(k8s.rawConfig)

	gvk := schema.FromAPIVersionAndKind("elbv2.k8s.aws/v1beta1", "TargetGroupBinding")
	gv := gvk.GroupVersion()
	target := gv.WithResource("targetgroupbindings")

	deploymentName := k8s.getDeploymentName(stackName, serviceName)
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
	deploymentName := k8s.getDeploymentName(stackName, serviceName)

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

// This function is used to retrieve events for a given stack, it looks into Deployment, Pod, Ingress, HorizontalPodAutoscaler and TargetGroupBinding triggered events
func (k8s *K8SComputeBackend) GetEvents(ctx context.Context, stackName string, services []string) error {
	if len(services) == 0 {
		return nil
	}

	eventsFound := false

	for _, serviceName := range services {
		resourceEvents := make([]corev1.Event, 0)
		deploymentName := k8s.getDeploymentName(stackName, serviceName)

		pods, err := k8s.getPods(ctx, deploymentName)
		if err != nil {
			return errors.Wrap(err, "unable to retrieve a list of pods")
		}
		if len(pods.Items) == 0 {
			return errors.New("No matching pods found, unable to retrieve events")
		}

		// Get events for all pods in a deployment
		for _, pod := range pods.Items {
			events, err := k8s.getResourceEvents(ctx, pod.Name, "Pod")
			if err != nil {
				return errors.Wrap(err, "unable to retrieve events for a pod")
			}
			resourceEvents = append(resourceEvents, events.Items...)
		}

		// Get events for the deployment, skipping ReplicaSet events for now
		events, err := k8s.getResourceEvents(ctx, deploymentName, "Deployment")
		if err != nil {
			return errors.Wrap(err, "unable to retrieve events for a deployment")
		}
		resourceEvents = append(resourceEvents, events.Items...)

		// Get events for the deployment, skipping ReplicaSet events for now
		events, err = k8s.getResourceEvents(ctx, deploymentName, "Service")
		if err != nil {
			return errors.Wrap(err, "unable to retrieve events for a service")
		}
		resourceEvents = append(resourceEvents, events.Items...)

		// Get events for the horizontal pod autoscaler
		events, err = k8s.getResourceEvents(ctx, deploymentName, "HorizontalPodAutoscaler")
		if err != nil {
			return errors.Wrap(err, "unable to retrieve events for a horizontal pod autoscaler")
		}
		resourceEvents = append(resourceEvents, events.Items...)

		// Get events for the ingress
		events, err = k8s.getResourceEvents(ctx, deploymentName, "Ingress")
		if err != nil {
			return errors.Wrap(err, "unable to retrieve events for an ingress resource")
		}
		resourceEvents = append(resourceEvents, events.Items...)

		// Find all matching target group bindings, and events for them. Target groups are created from the ingress resource. Target
		// groups are labeled with the "ingress.k8s.aws/stack=service-<STACK_NAME>-<SERVICE_NAME>" label.
		targetGroupBindings, err := k8s.getTargetGroupBindings(ctx, stackName, serviceName)
		if err != nil {
			return errors.Wrap(err, "unable to retrieve a list of ALB target group bindings")
		}

		// Get events for all target group bindings
		for _, targetGroupBinding := range targetGroupBindings.Items {
			events, err = k8s.getResourceEvents(ctx, targetGroupBinding.GetName(), "TargetGroupBinding")
			if err != nil {
				return errors.Wrap(err, "unable to retrieve events for a target group binding")
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

		warningCount := 0
		for _, e := range resourceEvents {
			if e.Type == Warning {
				logrus.Warnf("%s/%s - %s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message)
				warningCount++
			} else {
				logrus.Infof("%s/%s - %s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message)
			}
		}

		if warningCount > 1 {
			logrus.Println()
			logrus.Println("Many \"Warning\" events - please check to see whether your service is crashing:")
			logrus.Infof("  happy --env %s logs %s %s", k8s.Backend.Conf().GetEnv(), stackName, serviceName)
		}
		eventsFound = eventsFound || len(resourceEvents) > 0
	}

	if !eventsFound {
		logrus.Info("No events found for this stack")
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
				logrus.Errorf("unable to retrieve a list of resources %s/%s in namespace %s: %s", resource.Kind, resource.Name, k8s.KubeConfig.Namespace, err.Error())
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
		//"involvedObject.kind": resourceKind,
		//"type": Warning,
	})

	events, err := k8s.ClientSet.CoreV1().Events(k8s.KubeConfig.Namespace).List(ctx, v1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve events for resource %s/%s", resourceKind, resourceName)
	}
	return events, nil
}
