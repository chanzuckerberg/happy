package util

import "os"

func IsLocalstack() bool {
	mode, ok := os.LookupEnv("HAPPY_MODE")
	if !ok {
		return false
	}
	return mode == "localstack"
}
