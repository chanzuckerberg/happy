package util

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetStartTime gets the time the user started this command
func GetStartTime(ctx context.Context) time.Time {
	// This is the value the task was started, we don't want logs before this
	// time.
	cmdStartTime, ok := ctx.Value(CmdStartContextKey).(time.Time)
	if !ok {
		log.Debugf("didn't get a cmd start time. using now")
		cmdStartTime = time.Now()
	}
	return cmdStartTime
}
