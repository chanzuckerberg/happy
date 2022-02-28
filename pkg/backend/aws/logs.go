package aws

import (
	"context"

	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type GetLogsFunc func(*cwlv2.GetLogEventsOutput, error) error

func (b *Backend) getLogs(
	ctx context.Context,
	input *cwlv2.GetLogEventsInput,
	f GetLogsFunc,
) error {
	paginator := cwlv2.NewGetLogEventsPaginator(
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
