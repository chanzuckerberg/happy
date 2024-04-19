package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	str := "test"
	result := String(str)
	assert.NotNil(t, result)
	assert.Equal(t, str, *result)
}
