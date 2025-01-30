package log

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Formatter implements logrus.Formatter interface.
type Formatter struct{}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// This is inneficient but whatever
	var levelColor color.Attribute
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = color.FgBlack
	case logrus.WarnLevel:
		levelColor = color.FgYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.FgRed
	case logrus.InfoLevel:
		levelColor = color.FgBlue
	default:
		levelColor = color.FgHiBlack
	}

	level := color.New(levelColor).Sprintf("%s", strings.ToUpper(entry.Level.String()))

	messageColorizer := color.New(color.FgHiBlack).Sprintf

	message := bytes.NewBuffer(nil)
	_, err := message.WriteString(fmt.Sprintf("[%s]: %s\n", level, messageColorizer(entry.Message)))
	if err != nil {
		return nil, errors.Wrap(err, "could not compose message")
	}

	for key, value := range entry.Data {
		var toPrint interface{}
		switch v := value.(type) {
		case []byte:
			toPrint = string(v)
		default:
			toPrint = v
		}

		_, err = message.WriteString(messageColorizer("\t * %s: %+v\n", key, toPrint))
		if err != nil {
			return nil, err
		}
	}

	return message.Bytes(), nil
}
