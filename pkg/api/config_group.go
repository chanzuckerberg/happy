package api

import (
	"regexp"

	"github.com/chanzuckerberg/happy-api/pkg/cmd/config"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/chanzuckerberg/happy-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func RegisterConfigV1(v1 *fiber.Router) {
	group := (*v1).Group("/config")
	group.Get("/health", request.HealthHandler)

	// debugging endpoint that returns all config values for an app+env combo without resolving
	group.Get("/dump", parsePayload[model.AppMetadata], configDumpHandler)

	group.Post("/copy", parsePayload[model.CopyAppConfigPayload], configCopyHandler)
	group.Get("/diff", parsePayload[model.AppConfigDiffPayload], configDiffHandler)
	group.Post("/copyDiff", parsePayload[model.AppConfigDiffPayload], configCopyDiffHandler)

	loadConfigs(v1)
}

func loadConfigs(v1 *fiber.Router) {
	group := (*v1).Group("/configs")
	group.Get("/", parsePayload[model.AppMetadata], getConfigsHandler)
	group.Post("/", parsePayload[model.AppConfigPayload], postConfigsHandler)
	group.Get("/:key", parsePayload[model.AppMetadata], getConfigByKeyHandler)
	group.Delete("/:key", parsePayload[model.AppMetadata], deleteConfigByKeyHandler)
}

func getPayload[T interface{}](c *fiber.Ctx) T {
	return c.Context().UserValue("payload").(T)
}

func parsePayload[T interface{}](c *fiber.Ctx) error {
	payload := new(T)
	errors := request.ParsePayload(c, payload)
	if errors != nil {
		return response.ValidationErrorResponse(c, errors)
	}
	c.Context().SetUserValue("payload", *payload)

	return c.Next()
}

// @Summary Fetch all configs for a given app/env (including all stacks) (stack overrides are not applied)
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and optional stack"
// @Produce json
// @Success 200 {object} WrappedAppConfigsWithCount
// @Router  /v1/config/dump [GET]
func configDumpHandler(c *fiber.Ctx) error {
	payload := getPayload[model.AppMetadata](c)
	records, err := config.GetAllAppConfigs(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(wrapAppConfigsWithCount(records))
}

// @Summary Copy a single config key/value from one env/stack to another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, destination env/stack, and key to copy"
// @Produce json
// @Success 200 {object} WrappedResolvedAppConfig
// @Router  /v1/config/copy [POST]
func configCopyHandler(c *fiber.Ctx) error {
	payload := getPayload[model.CopyAppConfigPayload](c)

	record, err := config.CopyAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(WrappedAppConfig{Record: record})
}

// @Summary Get a list of config keys that exist in one env/stack and not another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, and destination env/stack"
// @Produce json
// @Success 200 {object} response.ConfigDiffResponse
// @Router  /v1/config/diff [GET]
func configDiffHandler(c *fiber.Ctx) error {
	payload := getPayload[model.AppConfigDiffPayload](c)

	missingKeys, err := config.AppConfigDiff(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.ConfigDiffResponse{MissingKeys: missingKeys})
}

// @Summary Copy the missing configs from one env/stack to another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, and destination env/stack"
// @Produce json
// @Success 200 {object} WrappedAppConfigsWithCount
// @Router  /v1/config/copyDiff [POST]
func configCopyDiffHandler(c *fiber.Ctx) error {
	payload := getPayload[model.AppConfigDiffPayload](c)

	records, err := config.CopyAppConfigDiff(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(wrapAppConfigsWithCount(records))
}

// @Summary Retrieve resolved configs for the given app/env/stack
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Produce json
// @Success 200 {object} WrappedResolvedAppConfigsWithCount
// @Router  /v1/configs/ [GET]
func getConfigsHandler(c *fiber.Ctx) error {
	payload := getPayload[model.AppMetadata](c)

	var records []*model.ResolvedAppConfig
	var err error
	if payload.Stack == "" {
		records, err = config.GetAppConfigsForEnv(&payload)
	} else {
		records, err = config.GetAppConfigsForStack(&payload)
	}
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(wrapResolvedAppConfigsWithCount(records))
}

// @Summary Retrieve resolved configs for the given app/env/stack
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigPayload true "Specification of the app, env, stack (optional), app config key, and app config value"
// @Produce json
// @Success 200 {object} WrappedAppConfig
// @Router  /v1/configs/ [POST]
func postConfigsHandler(c *fiber.Ctx) error {
	payload := getPayload[model.AppConfigPayload](c)
	payload.Key = standardizeKey(payload.Key)
	record, err := config.SetConfigValue(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(WrappedAppConfig{Record: record})
}

// @Summary Retrieve a single resolved config for the given app/env/stack and key
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Param   key     path string            true "The app config key to retrieve"
// @Produce json
// @Success 200 {object} WrappedResolvedAppConfig
// @Router  /v1/configs/{key} [GET]
func getConfigByKeyHandler(c *fiber.Ctx) error {
	payload := model.AppConfigLookupPayload{
		AppMetadata: getPayload[model.AppMetadata](c),
		ConfigKey:   model.ConfigKey{Key: c.Params("key")},
	}
	record, err := config.GetResolvedAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	status := c.Status(fiber.StatusOK)
	if record == nil {
		status = c.Status(fiber.StatusNotFound)
	}

	return status.JSON(WrappedResolvedAppConfig{Record: record})
}

// @Summary Delete a single resolved config for the given app/env/stack and key
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Param   key     path string            true "The app config key to delete"
// @Produce json
// @Success 200 {object} WrappedResolvedAppConfigsWithCount "record will be the deleted record (or null if nothing was deleted)"
// @Failure 400 {object} response.ValidationError
// @Router  /v1/configs/{key} [DELETE]
func deleteConfigByKeyHandler(c *fiber.Ctx) error {
	payload := model.AppConfigLookupPayload{
		AppMetadata: getPayload[model.AppMetadata](c),
		ConfigKey:   model.ConfigKey{Key: c.Params("key")},
	}
	record, err := config.DeleteAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(WrappedAppConfig{Record: record})
}

type WrappedAppConfigsWithCount struct {
	Records []*model.AppConfig `json:"records"`
	Count   int                `json:"count" example:"1"`
} // @Name response.WrappedAppConfigsWithCount

type WrappedResolvedAppConfigsWithCount struct {
	Records []*model.ResolvedAppConfig `json:"records"`
	Count   int                        `json:"count" example:"1"`
} // @Name response.WrappedResolvedAppConfigsWithCount

func wrapResolvedAppConfigsWithCount(records []*model.ResolvedAppConfig) WrappedResolvedAppConfigsWithCount {
	return WrappedResolvedAppConfigsWithCount{
		Records: records,
		Count:   len(records),
	}
}

func wrapAppConfigsWithCount(records []*model.AppConfig) WrappedAppConfigsWithCount {
	return WrappedAppConfigsWithCount{
		Records: records,
		Count:   len(records),
	}
}

// @Description App config key/value pair wrapped in "record" key
type WrappedResolvedAppConfig = struct {
	Record *model.ResolvedAppConfig `json:"record"`
} // @Name response.WrappedResolvedAppConfig

// @Description App config key/value pair wrapped in "record" key
type WrappedAppConfig = struct {
	Record *model.AppConfig `json:"record"`
} // @Name response.WrappedAppConfig

func standardizeKey(key string) string {
	// replace all non-alphanumeric characters with _
	regex := regexp.MustCompile("[^a-zA-Z0-9]")
	return regex.ReplaceAllString(key, "_")
}
