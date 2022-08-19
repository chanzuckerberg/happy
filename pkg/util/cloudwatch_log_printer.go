package util

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/fatih/color"
)

type logTemplate func(types.FilteredLogEvent) string

func timeStampeStreamMessageTemplate(event types.FilteredLogEvent) string {
	return fmt.Sprintf("[%.13d][%.15s] ", *event.Timestamp, *event.LogStreamName)
}

type TemplateOption func(lp *LogPrinter)
type LogPrinter struct {
	writer         io.Writer
	applyTemplate  logTemplate
	colors         []color.Attribute
	selectedColors map[string]color.Attribute
}

func WithLogTemplate(template logTemplate) TemplateOption {
	return func(lp *LogPrinter) {
		lp.applyTemplate = template
	}
}

func WithColors(colors []color.Attribute) TemplateOption {
	return func(lp *LogPrinter) {
		lp.colors = colors
	}
}

func WithWriter(writer io.Writer) TemplateOption {
	return func(lp *LogPrinter) {
		lp.writer = writer
	}
}

func MakeLogPrinter(opts ...TemplateOption) *LogPrinter {
	lp := LogPrinter{
		writer:         os.Stdout,
		selectedColors: map[string]color.Attribute{},
		applyTemplate:  timeStampeStreamMessageTemplate,
		colors: []color.Attribute{
			color.FgBlue,
			color.FgGreen,
			color.FgCyan,
			color.FgMagenta,
			color.FgYellow,
		},
	}
	for _, opt := range opts {
		opt(&lp)
	}

	return &lp
}

func (lp *LogPrinter) Log(fleo *cloudwatchlogs.FilterLogEventsOutput, err error) error {
	if err != nil {
		return err
	}

	for _, event := range fleo.Events {
		attr, ok := lp.selectedColors[*event.LogStreamName]
		if !ok {
			// no more colors to pick
			if len(lp.selectedColors) == len(lp.colors) {
				attr = color.FgBlack
			} else {
				attr = lp.colors[len(lp.selectedColors)]
			}
			lp.selectedColors[*event.LogStreamName] = attr
		}
		lp.writer.Write([]byte(color.New(attr).Sprintf("%s", lp.applyTemplate(event))))
		lp.writer.Write([]byte(fmt.Sprintf("%s\n", *event.Message)))
	}
	return nil
}
