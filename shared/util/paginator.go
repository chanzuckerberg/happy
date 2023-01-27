package util

import (
	"bufio"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Paginator interface {
	About()
	HasMorePages() bool
	NextPage(ctx context.Context) ([]LogEvent, error)
	WithSince(since int64) Paginator
	Build(ctx context.Context) (Paginator, error)
	Close(ctx context.Context) error
}

type CloudWatchPaginator struct {
	input         cloudwatchlogs.FilterLogEventsInput
	paginatorImpl *cloudwatchlogs.FilterLogEventsPaginator
}

type PodLogPaginator struct {
	podName    string
	logs       io.ReadCloser
	scanner    *bufio.Scanner
	pod        v1.PodInterface
	logOptions corev1.PodLogOptions
}

func NewCloudWatchPaginator(input cloudwatchlogs.FilterLogEventsInput, client cloudwatchlogs.FilterLogEventsAPIClient) Paginator {
	return &CloudWatchPaginator{
		input:         input,
		paginatorImpl: cloudwatchlogs.NewFilterLogEventsPaginator(client, &input),
	}
}

func NewPodLogPaginator(podName string, pod v1.PodInterface, logOptions corev1.PodLogOptions) Paginator {
	return &PodLogPaginator{
		podName:    podName,
		pod:        pod,
		logOptions: logOptions,
	}
}

func (p *CloudWatchPaginator) About() {
	logrus.Debugf("printing log group: '%s', log stream: '%+v'", *p.input.LogGroupName, p.input.LogStreamNames)
}

func (p *CloudWatchPaginator) HasMorePages() bool {
	return p.paginatorImpl.HasMorePages()
}

func (p *CloudWatchPaginator) NextPage(ctx context.Context) ([]LogEvent, error) {
	events := []LogEvent{}
	cloudwatchEvents, err := p.paginatorImpl.NextPage(ctx)
	if err != nil {
		return events, err
	}

	for _, event := range cloudwatchEvents.Events {
		events = append(events, LogEvent{
			Timestamp:     *event.Timestamp,
			LogStreamName: *event.LogStreamName,
			Message:       *event.Message,
		})
	}
	return events, nil
}

func (p *CloudWatchPaginator) WithSince(since int64) Paginator {
	p.input.StartTime = &since
	return p
}

func (p *CloudWatchPaginator) Build(ctx context.Context) (Paginator, error) {
	return p, nil
}

func (p *CloudWatchPaginator) Close(ctx context.Context) error {
	return nil
}

func (p *PodLogPaginator) About() {
	logrus.Debugf("printing logs from '%s'", p.podName)
}

func (p *PodLogPaginator) HasMorePages() bool {
	return p.scanner.Scan()
}

func (p *PodLogPaginator) NextPage(ctx context.Context) ([]LogEvent, error) {
	events := []LogEvent{
		{
			Timestamp:     0,
			LogStreamName: p.podName,
			Message:       string(p.scanner.Bytes()),
		},
	}
	return events, nil
}

func (p *PodLogPaginator) WithSince(since int64) Paginator {
	seconds := since / int64(time.Microsecond)
	nanoSeconds := (since % int64(time.Microsecond)) * int64(time.Millisecond)

	p.logOptions.SinceTime = &metav1.Time{
		Time: time.Unix(seconds, nanoSeconds),
	}
	return p
}

func (p *PodLogPaginator) Build(ctx context.Context) (Paginator, error) {
	logs, err := p.pod.GetLogs(p.podName, &p.logOptions).Stream(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve logs from pod %s", p.podName)
	}
	p.logs = logs
	p.scanner = bufio.NewScanner(logs)
	return p, nil
}

func (p *PodLogPaginator) Close(ctx context.Context) error {
	return p.logs.Close()
}
