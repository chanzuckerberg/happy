package cmd

import (
	"net/http"

	"github.com/chanzuckerberg/happy/pkg/cli/config"
	"github.com/chanzuckerberg/happy/pkg/cli/util"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ValidateWithHappyApi(cmd *cobra.Command, happyConfig *config.HappyConfig) error {
	if happyConfig.GetHappyApiBaseUrl() == "" {
		return errors.Errorf("Cannot use the %s feature set until you specify a valid happy-api URL in your happy config json", cmd.Use)
	}
	resp, err := util.MakeApiClient(happyConfig).Get("/versionCheck", nil)
	if err != nil {
		return errors.Wrap(err, "failed client version check")
	}

	if resp.StatusCode != http.StatusOK {
		jsonBody := map[string]interface{}{}
		err := client.ParseResponse(resp, &jsonBody)
		if err != nil {
			return err
		}
		return errors.Errorf("user-agent header used to validate your happy CLI version resulted in error: %s", jsonBody["message"])
	}
	return nil
}
