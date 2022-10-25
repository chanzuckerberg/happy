package util

import (
	"bufio"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
)

type Paginator interface {
	About()
	HasMorePages() bool
	NextPage(ctx context.Context) ([]LogEvent, error)
	WithSince(since int64)
}

type CloudWatchPaginator struct {
	input         cloudwatchlogs.FilterLogEventsInput
	paginatorImpl *cloudwatchlogs.FilterLogEventsPaginator
}

type ReaderPaginator struct {
	since   int64
	source  string
	scanner *bufio.Scanner
}

func NewCloudWatchPaginator(input cloudwatchlogs.FilterLogEventsInput, client cloudwatchlogs.FilterLogEventsAPIClient) Paginator {
	return &CloudWatchPaginator{
		input:         input,
		paginatorImpl: cloudwatchlogs.NewFilterLogEventsPaginator(client, &input),
	}
}

func NewReaderPaginator(source string, reader io.ReadCloser) Paginator {
	return &ReaderPaginator{
		source:  source,
		scanner: bufio.NewScanner(reader),
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

func (p *CloudWatchPaginator) WithSince(since int64) {
	p.input.StartTime = &since
}

func (p *ReaderPaginator) About() {
	logrus.Debugf("printing logs from '%s'", p.source)
}

func (p *ReaderPaginator) HasMorePages() bool {
	return p.scanner.Scan()
}

func (p *ReaderPaginator) NextPage(ctx context.Context) ([]LogEvent, error) {
	events := []LogEvent{
		{
			Timestamp:     0,
			LogStreamName: p.source,
			Message:       string(p.scanner.Bytes()),
		},
	}
	return events, nil
}

func (p *ReaderPaginator) WithSince(since int64) {
	p.since = since
}
