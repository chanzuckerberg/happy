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

	var attempt int64
	startTime := time.Now()
	for {
		out, err := b.cwlGetLogEventsAPIClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        input.LogGroupName,
			LogStreamNamePrefix: input.LogStreamName,
			Descending:          aws.Bool(true),
			Limit:               aws.Int32(10),
		})
		if err != nil {
			log.Errorf("error describing log streams: %s, retrying.", err.Error())
		} else {
			if len(out.LogStreams) > 0 {
				found := false
				for _, stream := range out.LogStreams {
					if *stream.LogStreamName == *input.LogStreamName {
						log.Infof("found a log stream after %d attempts", attempt)
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if time.Since(startTime) > 5*time.Minute {
			return errors.New("timed out waiting for log streams to be created")
		}
		if attempt > 100 {
			return errors.New("exceeded maximum number of attempts waiting for log streams to be created")
		}
		attempt++
		time.Sleep(time.Second * time.Duration(attempt*attempt))
	}

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
