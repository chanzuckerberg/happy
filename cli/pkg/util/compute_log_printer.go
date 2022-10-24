package util

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type logTemplate func(types.FilteredLogEvent) string

func timeStampeStreamMessageTemplate(event types.FilteredLogEvent) string {
	return fmt.Sprintf("[%.13d][%.15s] ", *event.Timestamp, *event.LogStreamName)
}

type PrintOption func(lp *ComputeLogPrinter)
type ComputeLogPrinter struct {
	writer         io.Writer
	input          cloudwatchlogs.FilterLogEventsInput
	applyTemplate  logTemplate
	colors         []color.Attribute
	selectedColors map[string]color.Attribute
}

func WithLogTemplate(template logTemplate) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.applyTemplate = template
	}
}

func WithColors(colors []color.Attribute) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.colors = colors
	}
}

func WithWriter(writer io.Writer) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.writer = writer
	}
}

func WithSince(since int64) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.input.StartTime = &since
	}
}

func WithCloudwatchInput(input cloudwatchlogs.FilterLogEventsInput) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.input = input
	}
}

func MakeComputeLogPrinter(opts ...PrintOption) *ComputeLogPrinter {
	lp := ComputeLogPrinter{
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

func (lp *ComputeLogPrinter) log(fleo *cloudwatchlogs.FilterLogEventsOutput, err error) error {
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
		logLine := fmt.Sprintf("%s%s", color.New(attr).Sprintf("%s", lp.applyTemplate(event)), fmt.Sprintf("%s\n", *event.Message))
		_, err = lp.writer.Write([]byte(logLine))
		if err != nil {
			return errors.Wrap(err, "error writing cloudwatch log")
		}
	}
	return nil
}

func (lp *ComputeLogPrinter) Print(ctx context.Context, client cloudwatchlogs.FilterLogEventsAPIClient) error {
	logrus.Debugf("printing log group: '%s', log stream: '%+v'", *lp.input.LogGroupName, lp.input.LogStreamNames)
	defer func() {
		logrus.Debug("cloudwatch log stream ended")
	}()

	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(client, &lp.input)
	for paginator.HasMorePages() {
		err := lp.log(paginator.NextPage(ctx))
		if IsStop(err) {
			return nil
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (lp *ComputeLogPrinter) PrintReader(ctx context.Context, source string, reader io.ReadCloser) error {
	logrus.Debug("printing logs")
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		attr, ok := lp.selectedColors[source]
		if !ok {
			// no more colors to pick
			if len(lp.selectedColors) == len(lp.colors) {
				attr = color.FgBlack
			} else {
				attr = lp.colors[len(lp.selectedColors)]
			}
			lp.selectedColors[source] = attr
		}
		logLine := color.New(attr).Sprintf("%s\n", string(scanner.Bytes()))
		_, err := lp.writer.Write([]byte(logLine))
		if IsStop(err) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "error writing cloudwatch log")
		}
	}
	return nil
}

var errStop = errors.New("stop")

func IsStop(err error) bool {
	return errors.Is(err, errStop)
}

// Pagination consumers can emit a Stop() to stop pagination
func Stop() error {
	return errStop
}
