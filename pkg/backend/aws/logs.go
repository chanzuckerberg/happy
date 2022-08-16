package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GetLogsFunc func(*cloudwatchlogs.GetLogEventsOutput, error) error

func (b *Backend) GetLogs(
	ctx context.Context,
	input *cloudwatchlogs.GetLogEventsInput,
	f GetLogsFunc,
) error {
	// ECS tasks produce per-container log streams, that only appear after a brief delay. Since timing is not guaranteed, we will re-try
	// with an exponential backoff.
	log.Infof("Waiting for the cloudwatch log stream to appear. Log group: '%s', log stream: '%s'", *input.LogGroupName, *input.LogStreamName)

	intervalWithTimeout(func() (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
		out, err := b.cwlGetLogEventsAPIClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        input.LogGroupName,
			LogStreamNamePrefix: input.LogStreamName,
			Descending:          aws.Bool(true),
			Limit:               aws.Int32(10),
		})
		if err != nil {
			log.Errorf("error describing log streams: %s, retrying.", err.Error())
			return nil, err
		}

		if len(out.LogStreams) > 0 {
			for _, stream := range out.LogStreams {
				if *stream.LogStreamName == *input.LogStreamName {
					return out, nil
				}
			}
		}
		return nil, errors.New("unable to find the log stream")
	}, 1*time.Second, 5*time.Minute)

	log.Info("\n...streaming cloudwatch logs...")

	// To tail, but not follow the log, specify options.StopOnDuplicateToken = true
	paginator := cloudwatchlogs.NewGetLogEventsPaginator(
		b.cwlGetLogEventsAPIClient,
		input,
	)

	for paginator.HasMorePages() {
		err := f(paginator.NextPage(ctx))
		if isStop(err) {
			return nil
		}

		if err != nil {
			log.Infof("error getting cloudwatch logs: %s", err.Error())
			return err
		}
	}
	log.Info("...cloudwatch log stream ended...")
	return nil
}
