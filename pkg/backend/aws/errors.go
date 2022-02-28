package aws

import (
	"github.com/pkg/errors"
)

var stop = errors.New("stop")

func isStop(err error) bool {
	return errors.Is(err, stop)
}

// Pagination consumers can emit a Stop() to stop pagination
func Stop() error {
	return stop
}
