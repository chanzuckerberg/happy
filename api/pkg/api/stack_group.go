package api

import (
	"fmt"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/gofiber/fiber/v2"
)

type StackHandler struct {
	// stack cmd.Stack
	stacklist cmd.Stacklist
}

func MakeStackHandler(s cmd.Stacklist) *StackHandler {
	return &StackHandler{
		stacklist: s,
	}
}

func MakeStackHandler2() *StackHandler {
	return &StackHandler{}
}

func (s *StackHandler) stacklistTestHandler(ctx *fiber.Ctx) error {
	// payload := new(AppStackPayload2)
	// errors := request.ParsePayload(ctx, payload)
	// if errors != nil {
	// 	return response.ValidationErrorResponse(ctx, errors)
	// }

	payload := getPayload[model.AppStackPayload2](ctx)
	fmt.Println(payload)

	stacks, err := s.stacklist.GetAppStacks(ctx.Context(), payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(wrapAppStacksWithCount(stacks))
}

func RegisterStackListV1(v1 fiber.Router, baseHandler *StackHandler) {
	v1.Get("/stacklistTest", parsePayload[model.AppStackPayload2], baseHandler.stacklistTestHandler)

	// group := v1.Group("/stacklistItems")
	// group.Get("/", parsePayload[model.AppStackPayload], baseHandler.getAppStacksHandler)
	// group.Post("/", parsePayload[model.AppStackPayload], baseHandler.createOrUpdateAppStackHandler)
	// group.Delete("/", parsePayload[model.AppStackPayload], baseHandler.deleteAppStackHandler)
}

// @Summary Retrieve app stacks for the given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app and env"
// @Produce json
// @Success 200 {object} model.WrappedAppStacksWithCount
// @Router  /v1/stacklistItems/ [GET]
// func (s *StackHandler) getAppStacksHandler(ctx *fiber.Ctx) error {
// 	payload := getPayload[model.AppStackPayload](ctx)
// 	stacks, err := s.stacklist.GetAppStacks(ctx.Context(), payload)
// 	if err != nil {
// 		return response.ServerErrorResponse(ctx, err.Error())
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(wrapAppStacksWithCount(stacks))
// }

// @Summary Create an app stack given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppStack true "Specification of the stack"
// @Produce json
// @Success 200 {object} model.WrappedAppStack
// @Router  /v1/stacklistItems/ [POST]
// func (s *StackHandler) createOrUpdateAppStackHandler(ctx *fiber.Ctx) error {
// 	payload := getPayload[model.AppStackPayload](ctx)
// 	stack, err := s.stack.CreateOrUpdateAppStack(payload)
// 	if err != nil {
// 		return response.ServerErrorResponse(ctx, err.Error())
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(wrapAppStack(stack))
// }

// @Summary Deletes a stack for the given app/env
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the stack"
// @Produce json
// @Success 200 {object} model.WrappedAppStack
// @Router  /v1/stacklistItems/ [DELETE]
// func (s *StackHandler) deleteAppStackHandler(ctx *fiber.Ctx) error {
// 	payload := getPayload[model.AppStackPayload](ctx)
// 	stack, err := s.stack.DeleteAppStack(payload)
// 	if err != nil {
// 		return response.ServerErrorResponse(ctx, err.Error())
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(wrapAppStack(stack))
// }

func wrapAppStacksWithCount(records []*model.AppStack) model.WrappedAppStacksWithCount {
	return model.WrappedAppStacksWithCount{
		Records: records,
		Count:   len(records),
	}
}

func wrapAppStack(record *model.AppStack) model.WrappedAppStack {
	return model.WrappedAppStack{
		Record: record,
	}
}
