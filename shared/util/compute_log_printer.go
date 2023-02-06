package util

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogEvent struct {
	LogStreamName string
	Message       string
	Timestamp     int64
}

type logTemplate func(LogEvent) string

func timeStampedStreamMessageTemplate(event LogEvent) string {
	return fmt.Sprintf("[%.13d][%.15s] ", event.Timestamp, event.LogStreamName)
}

func RawStreamMessageTemplate(event LogEvent) string {
	return ""
}

type PrintOption func(lp *ComputeLogPrinter)
type ComputeLogPrinter struct {
	writer         io.Writer
	paginator      Paginator
	applyTemplate  logTemplate
	colors         []color.Attribute
	selectedColors map[string]color.Attribute
	since          int64
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

func WithPaginator(paginator Paginator) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.paginator = paginator
	}
}

func WithSince(since int64) PrintOption {
	return func(lp *ComputeLogPrinter) {
		lp.since = since
	}
}

func MakeComputeLogPrinter(ctx context.Context, opts ...PrintOption) *ComputeLogPrinter {
	lp := ComputeLogPrinter{
		writer:         os.Stdout,
		selectedColors: map[string]color.Attribute{},
		applyTemplate:  timeStampedStreamMessageTemplate,
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
	lp.paginator.WithSince(lp.since)

	return &lp
}

func (lp *ComputeLogPrinter) log(events []LogEvent, err error) error {
	if err != nil {
		return err
	}
	for _, event := range events {
		attr, ok := lp.selectedColors[event.LogStreamName]
		if !ok {
			// no more colors to pick
			if len(lp.selectedColors) == len(lp.colors) {
				attr = color.FgBlack
			} else {
				attr = lp.colors[len(lp.selectedColors)]
			}
			lp.selectedColors[event.LogStreamName] = attr
		}
		logLine := fmt.Sprintf("%s%s", color.New(attr).Sprintf("%s", lp.applyTemplate(event)), fmt.Sprintf("%s\n", event.Message))
		_, err := lp.writer.Write([]byte(logLine))
		if err != nil {
			return errors.Wrap(err, "error writing log")
		}
	}
	return nil
}

func (lp *ComputeLogPrinter) Print(ctx context.Context) error {
	_, err := lp.paginator.Build(ctx)
	if err != nil {
		return errors.Wrap(err, "Cannot set up log pagination")
	}
	lp.paginator.About()
	defer func() {
		logrus.Debug("log stream ended")
	}()

	for lp.paginator.HasMorePages() {
		events, err := lp.paginator.NextPage(ctx)
		err = lp.log(events, err)

		if IsStop(err) {
			return nil
		}

		if err != nil {
			return err
		}
	}
	err = lp.paginator.Close(ctx)
	if err != nil {
		return errors.Wrap(err, "Cannot wrap up logging")
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
