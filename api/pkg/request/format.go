package request

import (
	"regexp"
	"strings"
)

func StandardizeKey(key string) string {
	key = strings.ToUpper(key)

	// replace all non-alphanumeric characters with _
	regex := regexp.MustCompile("[^A-Z0-9]")
	return regex.ReplaceAllString(key, "_")
}
