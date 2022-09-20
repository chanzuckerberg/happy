package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ValidateHappyConfig(cmd *cobra.Command, args []string) error {
	happyConfig, err := GetHappyConfig(cmd)
	if err != nil {
		return err
	}

	if !happyConfig.GetFeatures().EnableHappyApiUsage {
		return errors.Errorf("Cannot use the %s feature set until you enable happy-api usage in your happy config json", cmd.Use)
	}
	if happyConfig.GetHappyApiBaseUrl() == "" {
		return errors.Errorf("Cannot use the %s feature set until you specify a valid happy-api URL in your happy config json", cmd.Use)
	}
	resp, err := NewHappyClient(happyConfig).Get("/versionCheck", nil)
	if err != nil {
		return errors.Wrap(err, "failed client version check")
	}

	if resp.StatusCode != http.StatusOK {
		jsonBody := map[string]interface{}{}
		err := ParseResponse(resp, &jsonBody)
		if err != nil {
			return err
		}
		return errors.Errorf("user-agent header used to validate your happy CLI version resulted in error: %s", jsonBody["message"])
	}

	return nil
}

func GetHappyConfig(cmd *cobra.Command) (*config.HappyConfig, error) {
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	return happyConfig, nil
}

func ParseResponse[T interface{}](resp *http.Response, result *T) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	return nil
}
