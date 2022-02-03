package util

import (
	"encoding/json"
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

func TagValueToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch t := value.(type) {
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case string:
		return value.(string)
	case map[string]interface{}:
		if len(t) == 0 {
			return ""
		} else {
			data, err := json.Marshal(t)
			if err != nil {
				return ""
			}
			return string(data)
		}
	default:
		return fmt.Sprintf("%v", t)
	}
}
