package util

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntervalWithTimeout_Success(t *testing.T) {
	// Define a mock function that returns a value without error
	mockFunc := func() (int, error) {
		return 42, nil
	}

	// Call the IntervalWithTimeout function with a short timeout
	result, err := IntervalWithTimeout(mockFunc, time.Millisecond, time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 42, *result)
}

func TestIntervalWithTimeout_Timeout(t *testing.T) {
	// Define a mock function that always returns an error
	mockFunc := func() (int, error) {
		return 0, errors.New("mock error")
	}

	// Call the IntervalWithTimeout function with a very short timeout
	result, err := IntervalWithTimeout(mockFunc, time.Millisecond, time.Millisecond)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "timed out")
}
