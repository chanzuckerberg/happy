package util

import (
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

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
			return "{}"
		}
		data, err := json.Marshal(t)
		if err != nil {
			log.Debugf("Cannot serialize to json: %v\n", value)
			return ""
		}
		return string(data)
	default:
		return fmt.Sprintf("%v", t)
	}
}
