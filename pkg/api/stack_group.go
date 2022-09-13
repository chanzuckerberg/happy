package api

import (
	"github.com/chanzuckerberg/happy-api/pkg/cmd"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/chanzuckerberg/happy-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type StackHandler struct {
	stack cmd.Stack
}

func MakeStackHandler(s cmd.Stack) *StackHandler {
	return &StackHandler{
		stack: s,
	}
}

func RegisterStackListV1(v1 fiber.Router, baseHandler *StackHandler) {
	group := v1.Group("/stacks")
	group.Get("/", parsePayload[model.AppStackPayload], baseHandler.getAppStacksHandler)
	group.Post("/", parsePayload[model.AppStackPayload], baseHandler.createAppStackHandler)
	group.Put("/", parsePayload[model.AppStackPayload], baseHandler.updateAppStackHandler)
}

// @Summary Retrieve app stacks for the given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app and env"
// @Produce json
// @Success 200 {object} WrappedAppStacksWithCount
// @Router  /v1/stacks/ [GET]
func (s *StackHandler) getAppStacksHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppStackPayload](ctx)
	stacks, err := s.stack.GetAppStacks(payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppStacksWithCount(stacks))
}

// @Summary Create an app stack given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppStack true "Specification of the stack"
// @Produce json
// @Success 200 {object} WrappedAppStack
// @Router  /v1/stacks/ [POST]
func (s *StackHandler) createAppStackHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppStackPayload](ctx)
	stack, err := s.stack.CreateAppStack(payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppStack(stack))
}

// @Summary Updates the enabled column of a stack for the given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppStack true "Specification of the stack"
// @Produce json
// @Success 200 {object} WrappedAppStack
// @Router  /v1/stacks/ [PUT]
func (s *StackHandler) updateAppStackHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppStackPayload](ctx)
	stack, err := s.stack.UpdateAppStack(payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppStack(stack))
}

type WrappedAppStacksWithCount struct {
	Records []*model.AppStack `json:"records"`
	Count   int               `json:"count" example:"1"`
} // @Name response.WrappedAppStacksWithCount

type WrappedAppStack struct {
	Record *model.AppStack `json:"record"`
} // @Name response.WrappedAppStack

func wrapAppStacksWithCount(records []*model.AppStack) WrappedAppStacksWithCount {
	return WrappedAppStacksWithCount{
		Records: records,
		Count:   len(records),
	}
}

func wrapAppStack(record *model.AppStack) WrappedAppStack {
	return WrappedAppStack{
		Record: record,
	}
}
