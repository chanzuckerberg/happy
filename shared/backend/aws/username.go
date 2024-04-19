package aws

import (
	"regexp"
	"strings"
)

var userNameRegex = regexp.MustCompile(`[^\p{L}\p{Z}\p{N}_.:\/=+\-@]+`)

func cleanupUserName(username string) string {
	username = userNameRegex.ReplaceAllString(username, "-")
	username = strings.Trim(username, "-")
	return username
}
