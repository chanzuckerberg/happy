package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LogMessages struct {
	// loggroup/logstream : messages
	messages map[string][]string
}

func (lm *LogMessages) Print() {
	for name, messages := range lm.messages {
		log.Infof("\n\nLOG Events (%s):", name)
		log.Println("************************************")
		for _, m := range messages {
			log.Println(m)
		}
		log.Println("************************************")
	}
}

func (b *Backend) getLogs(ctx context.Context, input *cloudwatchlogs.GetLogEventsInput) (*LogMessages, error) {
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
	name := fmt.Sprintf("%s/%s", *input.LogGroupName, *input.LogStreamName)
	return &LogMessages{
		messages: map[string][]string{
			name: messages,
		},
	}, nil
}
