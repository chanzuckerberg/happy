package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	log "github.com/sirupsen/logrus"
)

type GetLogsFunc func(*cloudwatchlogs.GetLogEventsOutput, error) error
type FilterLogsFunc interface {
	Log(fleo *cloudwatchlogs.FilterLogEventsOutput, err error) error
}

func (b *Backend) GetLogs(ctx context.Context, input *cloudwatchlogs.FilterLogEventsInput, f FilterLogsFunc) error {
	log.Debugf("log group: '%s', log stream: '%+v'", *input.LogGroupName, input.LogStreamNames)

	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(
		b.cwlFilterLogEventsAPIClient,
		input,
	)

	for paginator.HasMorePages() {
		err := f.Log(paginator.NextPage(ctx))
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
