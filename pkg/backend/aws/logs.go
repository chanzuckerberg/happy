package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	log "github.com/sirupsen/logrus"
)

type GetLogsFunc func(*cloudwatchlogs.GetLogEventsOutput, error) error

func (b *Backend) GetLogs(
	ctx context.Context,
	input *cloudwatchlogs.GetLogEventsInput,
	f GetLogsFunc,
) error {
	log.Infof("cloudwatch log group: '%s', log stream: '%s'", *input.LogGroupName, *input.LogStreamName)
	log.Info("\n...streaming cloudwatch logs...")
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
			return err
		}
	}
	log.Info("...cloudwatch log stream ended...")
	return nil
}
