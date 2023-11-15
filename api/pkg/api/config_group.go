package api

import (
	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/gofiber/fiber/v2"
)

type ConfigHandler struct {
	config cmd.Config
}

func MakeConfigHandler(c cmd.Config) *ConfigHandler {
	return &ConfigHandler{
		config: c,
	}
}

func RegisterConfigV1(v1 fiber.Router, baseHandler *ConfigHandler) {
	group := v1.Group("/config")
	group.Get("/health", request.HealthHandler)

	// debugging endpoint that returns all config values for an app+env combo without resolving
	group.Get("/dump", parseQueryString[model.AppMetadata], baseHandler.configDumpHandler)
	group.Post("/copy", parseRequestBody[model.CopyAppConfigPayload], baseHandler.configCopyHandler)
	group.Get("/diff", parseQueryString[model.AppConfigDiffPayload], baseHandler.configDiffHandler)
	group.Post("/copyDiff", parseRequestBody[model.AppConfigDiffPayload], baseHandler.configCopyDiffHandler)

	loadConfigs(v1, baseHandler)
}

func loadConfigs(v1 fiber.Router, baseHandler *ConfigHandler) {
	group := v1.Group("/configs")
	group.Get("/", parseQueryString[model.AppMetadata], baseHandler.getConfigsHandler)
	group.Post("/", parseRequestBody[model.AppConfigPayload], baseHandler.postConfigsHandler)
	group.Get("/:key", parseQueryString[model.AppMetadata], baseHandler.getConfigByKeyHandler)
	group.Delete("/:key", parseRequestBody[model.AppMetadata], baseHandler.deleteConfigByKeyHandler)
}

func getPayload[T interface{}](c *fiber.Ctx) T {
	return c.Context().UserValue("payload").(T)
}

// not needed in prod but useful to keep around for debugging
// func requestInspector(c *fiber.Ctx) error {
// 	fmt.Println("----------------------------")
// 	fmt.Println("Request Inspection Data")
// 	fmt.Printf("- Request Method: %s\n", c.Route().Method)
// 	fmt.Printf("- Request Route:  %s\n", c.Route().Path)
// 	fmt.Printf("- Query String:   %s\n", c.Context().QueryArgs())
// 	fmt.Println("----------------------------")

// 	return c.Next()
// }

func parseRequestBody[T interface{}](c *fiber.Ctx) error {
	return parsePayload[T](c, c.BodyParser)
}

func parseQueryString[T interface{}](c *fiber.Ctx) error {
	return parsePayload[T](c, c.QueryParser)
}

func parsePayload[T interface{}](c *fiber.Ctx, fn request.RequestParser) error {
	payload := new(T)
	errors := request.ParsePayload(c, payload, fn)
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
// @Success 200 {object} model.WrappedAppConfigsWithCount
// @Router  /v1/config/dump [GET]
func (c *ConfigHandler) configDumpHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppMetadata](ctx)
	records, err := c.config.GetAllAppConfigs(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppConfigsWithCount(records))
}

// @Summary Copy a single config key/value from one env/stack to another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, destination env/stack, and key to copy"
// @Produce json
// @Success 200 {object} model.WrappedResolvedAppConfig
// @Router  /v1/config/copy [POST]
func (c *ConfigHandler) configCopyHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.CopyAppConfigPayload](ctx)
	payload.Key = request.StandardizeKey(payload.Key)

	record, err := c.config.CopyAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WrappedAppConfig{Record: record})
}

// @Summary Get a list of config keys that exist in one env/stack and not another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, and destination env/stack"
// @Produce json
// @Success 200 {object} model.ConfigDiffResponse
// @Router  /v1/config/diff [GET]
func (c *ConfigHandler) configDiffHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppConfigDiffPayload](ctx)

	missingKeys, err := c.config.AppConfigDiff(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ConfigDiffResponse{MissingKeys: missingKeys})
}

// @Summary Copy the missing configs from one env/stack to another
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigDiffPayload true "Specification of the app, source env/stack, and destination env/stack"
// @Produce json
// @Success 200 {object} model.WrappedAppConfigsWithCount
// @Router  /v1/config/copyDiff [POST]
func (c *ConfigHandler) configCopyDiffHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppConfigDiffPayload](ctx)

	records, err := c.config.CopyAppConfigDiff(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppConfigsWithCount(records))
}

// @Summary Retrieve resolved configs for the given app/env/stack
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Produce json
// @Success 200 {object} model.WrappedResolvedAppConfigsWithCount
// @Router  /v1/configs/ [GET]
func (c *ConfigHandler) getConfigsHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppMetadata](ctx)

	var records []*model.ResolvedAppConfig
	var err error
	if payload.Stack == "" {
		records, err = c.config.GetAppConfigsForEnv(&payload)
	} else {
		records, err = c.config.GetAppConfigsForStack(&payload)
	}
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapResolvedAppConfigsWithCount(records))
}

// @Summary Retrieve resolved configs for the given app/env/stack
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppConfigPayload true "Specification of the app, env, stack (optional), app config key, and app config value"
// @Produce json
// @Success 200 {object} model.WrappedAppConfig
// @Router  /v1/configs/ [POST]
func (c *ConfigHandler) postConfigsHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppConfigPayload](ctx)
	payload.Key = request.StandardizeKey(payload.Key)
	record, err := c.config.SetConfigValue(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WrappedAppConfig{Record: record})
}

// @Summary Retrieve a single resolved config for the given app/env/stack and key
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Param   key     path string            true "The app config key to retrieve"
// @Produce json
// @Success 200 {object} model.WrappedResolvedAppConfig
// @Router  /v1/configs/{key} [GET]
func (c *ConfigHandler) getConfigByKeyHandler(ctx *fiber.Ctx) error {
	payload := model.AppConfigLookupPayload{
		AppMetadata: getPayload[model.AppMetadata](ctx),
		ConfigKey:   model.ConfigKey{Key: request.StandardizeKey(ctx.Params("key"))},
	}
	record, err := c.config.GetResolvedAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	status := ctx.Status(fiber.StatusOK)
	if record == nil {
		status = ctx.Status(fiber.StatusNotFound)
	}

	return status.JSON(model.WrappedResolvedAppConfig{Record: record})
}

// @Summary Delete a single resolved config for the given app/env/stack and key
// @Tags    config
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app, env, and stack (optional)"
// @Param   key     path string            true "The app config key to delete"
// @Produce json
// @Success 200 {object} model.WrappedAppConfig "record will be the deleted record (or null if nothing was deleted)"
// @Failure 400 {object} model.ValidationError
// @Router  /v1/configs/{key} [DELETE]
func (c *ConfigHandler) deleteConfigByKeyHandler(ctx *fiber.Ctx) error {
	payload := model.AppConfigLookupPayload{
		AppMetadata: getPayload[model.AppMetadata](ctx),
		ConfigKey:   model.ConfigKey{Key: request.StandardizeKey(ctx.Params("key"))},
	}
	record, err := c.config.DeleteAppConfig(&payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WrappedAppConfig{Record: record})
}

func wrapResolvedAppConfigsWithCount(records []*model.ResolvedAppConfig) model.WrappedResolvedAppConfigsWithCount {
	return model.WrappedResolvedAppConfigsWithCount{
		Records: records,
		Count:   len(records),
	}
}

func wrapAppConfigsWithCount(records []*model.AppConfig) model.WrappedAppConfigsWithCount {
	return model.WrappedAppConfigsWithCount{
		Records: records,
		Count:   len(records),
	}
}
