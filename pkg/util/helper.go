package util

import (
	"fmt"
	"time"

	"github.com/chanzuckerberg/happy-deploy/pkg/backend"
	"github.com/chanzuckerberg/happy-deploy/pkg/config"
)

func GenerateTag(config config.HappyConfigIface) (string, error) {
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
