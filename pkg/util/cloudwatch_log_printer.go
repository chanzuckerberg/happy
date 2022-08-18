package util

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type logTemplate func(types.FilteredLogEvent) string

func timeStampeSteamMessageTemplate(event types.FilteredLogEvent) string {
	return fmt.Sprintf("[%.20s][%.10s] ", time.Unix(*event.Timestamp, 0), *event.LogStreamName)
}

type TemplateOption func(lp *LogPrinter)
type LogPrinter struct {
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

func MakeLogPrinter(opts ...TemplateOption) *LogPrinter {
	lp := LogPrinter{
		selectedColors: map[string]color.Attribute{},
		applyTemplate:  timeStampeSteamMessageTemplate,
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
		color.New(attr).Printf("%s", lp.applyTemplate(event))
		logrus.Infof("%s\n", *event.Message)
	}
	return nil
}
