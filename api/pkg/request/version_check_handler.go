package request

import (
	"strings"

	"github.com/blang/semver"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

var (
	MinimumVersions map[string]string
)

func init() {
	MinimumVersions = map[string]string{
		"happy-cli":      "0.90.0",
		"happy-provider": "0.52.0",
	}

	// panic if any of the above versions are invalid version strings
	for _, versionStr := range MinimumVersions {
		semver.MustParse(versionStr)
	}
}

func VersionCheckHandler(c *fiber.Ctx) error {
	userAgent := string(c.Request().Header.UserAgent())
	err := validateUserAgentVersion(userAgent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": ""})
}

func validateUserAgentVersion(userAgent string) error {
	userAgentParts := strings.Fields(userAgent)
	if len(userAgentParts) == 0 {
		return nil
	}
	clientAndVersion := userAgentParts[0]
	clientVersionParts := strings.Split(clientAndVersion, "/")
	client := clientVersionParts[0]
	minVersionStr, clientFound := MinimumVersions[client]
	if !clientFound {
		return nil
	}

	if len(clientVersionParts) < 2 {
		return errors.Errorf("expected version to be specified for %s in the User-Agent header (format: %s/<version>)", client, client)
	}

	versionStr := clientVersionParts[1]
	version, err := semver.Parse(versionStr)
	if err != nil {
		return err
	}

	minVersion, err := semver.Parse(minVersionStr)
	if err != nil {
		return err
	}

	if clientFound && version.LT(minVersion) {
		return errors.Errorf("please upgrade your %s client to at least v%s", client, minVersionStr)
	}

	return nil
}
