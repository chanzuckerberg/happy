package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LogMessages struct {
	messages []string
}

func (lm *LogMessages) Print() {
	// TODO compact the print format of log event (currently prints single
	// field per line) to single line per event
	log.Println("\n\nLOG Events:")
	log.Println("************************************")
	for _, m := range lm.messages {
		log.Println(m)
	}
	log.Println("************************************")
}

func (b *awsBackend) getLogs(ctx context.Context, input *cloudwatchlogs.GetLogEventsInput) ([]string, error) {
	// TODO(el): do we want paging here?
	out, err := b.logsclient.GetLogEventsWithContext(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "could not get logs")
	}

	messages := []string{}
	for _, event := range out.Events {
		if event == nil {
			continue
		}
		msg := event.Message
		if msg == nil {
			continue
		}
		messages = append(messages, *msg)
	}
	return messages, nil
}
