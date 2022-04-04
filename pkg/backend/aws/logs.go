package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type GetLogsFunc func(*cloudwatchlogs.GetLogEventsOutput, error) error

func (b *Backend) GetLogs(
	ctx context.Context,
	input *cloudwatchlogs.GetLogEventsInput,
	f GetLogsFunc,
) error {
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
	return nil
}
