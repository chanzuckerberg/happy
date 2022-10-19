package client

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/model"
)

func (c *HappyClient) ListConfigs(appName, environment, stack string) (model.WrappedResolvedAppConfigsWithCount, error) {
	body := model.NewAppMetadata(appName, environment, stack)
	result := model.WrappedResolvedAppConfigsWithCount{}
	err := c.GetParsed("/v1/configs", body, &result)
	return result, err
}

func (c *HappyClient) GetConfig(appName, environment, stack, key string) (model.WrappedResolvedAppConfig, error) {
	body := model.NewAppMetadata(appName, environment, stack)
	result := model.WrappedResolvedAppConfig{}
	err := c.GetParsed(fmt.Sprintf("/v1/configs/%s", key), body, &result)
	return result, err
}

func (c *HappyClient) SetConfig(appName, environment, stack, key, value string) (model.WrappedResolvedAppConfig, error) {
	body := model.NewAppConfigPayload(appName, environment, stack, key, value)
	result := model.WrappedResolvedAppConfig{}
	err := c.PostParsed("/v1/configs", body, &result)
	return result, err
}

func (c *HappyClient) DeleteConfig(appName, environment, stack, key string) (model.WrappedResolvedAppConfig, error) {
	body := model.NewAppConfigLookupPayload(appName, environment, stack, key)
	result := model.WrappedResolvedAppConfig{}
	err := c.DeleteParsed(fmt.Sprintf("/v1/configs/%s", key), body, &result)
	return result, err
}

func (c *HappyClient) CopyConfig(appName, srcEnv, srcStack, destEnv, destStack, key string) (model.WrappedResolvedAppConfig, error) {
	body := model.NewCopyAppConfigPayload(appName, srcEnv, srcStack, destEnv, destStack, key)
	result := model.WrappedResolvedAppConfig{}
	err := c.PostParsed("/v1/config/copy", body, &result)
	return result, err
}

func (c *HappyClient) GetMissingConfigKeys(appName, srcEnv, srcStack, destEnv, destStack string) (model.ConfigDiffResponse, error) {
	body := model.NewAppConfigDiffPayload(appName, srcEnv, srcStack, destEnv, destStack)
	result := model.ConfigDiffResponse{}
	err := c.GetParsed("/v1/config/diff", body, &result)
	return result, err
}
