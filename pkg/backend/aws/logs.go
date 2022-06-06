package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type GetLogsFunc func(*cloudwatchlogs.GetLogEventsOutput, error) error

func (b *Backend) GetLogs(
	ctx context.Context,
	input *cloudwatchlogs.GetLogEventsInput,
	f GetLogsFunc,
) error {
	log.Debugf("cloudwatch log group: '%s', log stream: '%s'", *input.LogGroupName, *input.LogStreamName)
	log.Info("\n...streaming cloudwatch logs...")
	paginator := cloudwatchlogs.NewGetLogEventsPaginator(
		b.cwlGetLogEventsAPIClient,
		input,
	)
	// NOTE[JH](CCIE-220): According to cloudwatch documentation, GetLogEvents only
	// allows 25 requets / second. This limiter is configured for a little less than
	// that, with an initial bucket size so we stop hitting the rate limit when stream logs.
	// https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/cloudwatch_limits_cwl.html
	limiter := rate.NewLimiter(rate.Limit(20), 1)
	for paginator.HasMorePages() {
		err := limiter.Wait(ctx)
		if err != nil {
			return errors.Wrap(err, "error waiting for GetLogEvents rate limit to fill back up")
		}
		err = f(paginator.NextPage(ctx))
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
