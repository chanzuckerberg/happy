package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
)

func GenerateTag(config config.HappyConfig) (string, error) {
	t := time.Now()
	ts := fmt.Sprintf("%02d%02d-%02d%02d%02d", t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	userIdBackend := backend.GetAwsSts(config)
	userName, err := userIdBackend.GetUserName()
	if err != nil {
		return "", err
	}
	tag := fmt.Sprintf("%s-%s", userName, ts)

	return tag, nil
}

func ToString(value interface{}) (res string) {
	switch value.(type) {
	case float64:
		res = strconv.FormatFloat(value.(float64), 'f', 0, 64)
	case string:
		res = value.(string)
	default:
		res = ""
	}
	return
}
