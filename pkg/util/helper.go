package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	log "github.com/sirupsen/logrus"
)

// NOTE(el): This is based off RFC3339 with some tweaks to make it a valid docker tag
const dockerRFC3339TimeFmt string = "2006-01-02T15-04-05"

func GenerateTag(config config.HappyConfig) (string, error) {
	userIdBackend := backend.GetAwsSts(config)
	userName, err := userIdBackend.GetUserName()
	if err != nil {
		return "", err
	}
	t := time.Now().UTC().Format(dockerRFC3339TimeFmt)
	tag := fmt.Sprintf("%s-%s", userName, t)

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
