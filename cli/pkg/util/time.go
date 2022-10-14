package util

import (
	"context"
	"time"

	"github.com/pkg/errors"
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

// intervalWithTimeout is a helper function to run a function many times with a given interval and a set timeout period
func IntervalWithTimeout[K any](f func() (K, error), tick time.Duration, timeout time.Duration) (*K, error) {
	timeoutChan := time.After(timeout)
	tickChan := time.NewTicker(tick)

	for {
		select {
		case <-timeoutChan:
			return nil, errors.New("timed out")
		case <-tickChan.C:
			out, err := f()
			if err == nil {
				return &out, nil
			}
			log.Debugf("trying again: %s", err)
		}
	}
}
