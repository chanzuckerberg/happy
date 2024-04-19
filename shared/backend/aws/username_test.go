package aws

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidUserName(t *testing.T) {
	r := require.New(t)
	username := cleanupUserName("user@domain.com")
	r.Equal("user@domain.com", username)

	username = cleanupUserName("user")
	r.Equal("user", username)

	username = cleanupUserName("github-helper[bot]")
	r.Equal("github-helper-bot", username)
}
