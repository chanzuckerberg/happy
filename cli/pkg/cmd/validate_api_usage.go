package cmd

import (
	"net/http"

	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ValidateWithHappyApi(cmd *cobra.Command, happyConfig *config.HappyConfig, backend *aws.Backend) error {
	if happyConfig.GetHappyAPIConfig().BaseUrl == "" {
		return errors.Errorf("Cannot use the %s feature set until you specify a valid happy-api URL in your happy config json", cmd.Use)
	}
	resp, err := hapi.MakeAPIClient(happyConfig, backend).Get("/versionCheck", nil)
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
