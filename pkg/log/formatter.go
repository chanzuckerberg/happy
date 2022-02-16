package log

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

	level := color.New(levelColor).Sprintf(strings.ToUpper(entry.Level.String()))
	message := color.New(color.FgHiBlack).Sprint(entry.Message)

	output := fmt.Sprintf("[%s]: %s\n", level, message)
	return []byte(output), nil
}
