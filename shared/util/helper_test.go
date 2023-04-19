package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringTagValue(t *testing.T) {
	assert.Equal(t, "", TagValueToString(""), "empty string value should not change")
	assert.Equal(t, "value", TagValueToString("value"), "non-empty string value should not change")
}

func TestFloatTagValue(t *testing.T) {
	assert.Equal(t, "120.01", TagValueToString(120.01), "float value should be represented correctly")
	assert.Equal(t, "120", TagValueToString(120.00), "float value should be represented correctly")
	assert.Equal(t, "-1", TagValueToString(-1.0), "non-empty string value should not change")
	assert.Equal(t, "0.01", TagValueToString(0.01), "non-empty string value should not change")
}

func TestMapTagValue(t *testing.T) {
	assert.Equal(t, "", TagValueToString(map[string]interface{}{}), "empty map value should become blank")
	assert.Equal(t, "{\"key\":\"value\"}", TagValueToString(map[string]interface{}{"key": "value"}), "non-empty map value should become a valid json")
}

func TestOtherTagValue(t *testing.T) {
	assert.Equal(t, "", TagValueToString(nil), "nil value should become blank")
	assert.Equal(t, "1", TagValueToString(1), "int value should remain")
	assert.Equal(t, "-1", TagValueToString(-1), "negative int value should remain")
}
