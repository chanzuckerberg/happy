package api

import (
	"context"
	"fmt"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/gofiber/fiber/v2"
)

type StackHandler struct {
	stack cmd.StackManager
}

func MakeStackHandler(s cmd.StackManager) *StackHandler {
	return &StackHandler{
		stack: s,
	}
}

func RegisterStackListV1(v1 fiber.Router, baseHandler *StackHandler) {
	group := v1.Group("/stacks")
	group.Get("/", parseQueryString[model.AppStackPayload], baseHandler.getAppStacksHandler)
	// group.Post("/", parsePayload[model.AppStackPayload], baseHandler.createOrUpdateAppStackHandler)
	// group.Delete("/", parsePayload[model.AppStackPayload], baseHandler.deleteAppStackHandler)
}

// @Summary Retrieve app stacks for the given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppMetadata true "Specification of the app and env"
// @Produce json
// @Success 200 {object} model.WrappedAppStacksWithCount
// @Router  /v1/stacks/ [GET]
func (s *StackHandler) getAppStacksHandler(ctx *fiber.Ctx) error {
	payload := getPayload[model.AppStackPayload](ctx)
	// TODO: make to middleware so it can be reused
	authdCtx, err := request.CtxWithAWSAuthHeaders(ctx)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}

	stacks, err := s.stack.GetAppStacks(authdCtx, payload)
	if err != nil {
		return response.ServerErrorResponse(ctx, err.Error())
	}
	go func(c context.Context) {
		fmt.Println("going!")
		for {
			select {
			case <-c.Done():
				fmt.Println("done")
				return
			}
		}
	}(ctx.UserContext())

	return ctx.Status(fiber.StatusOK).JSON(wrapAppStacksWithCount(stacks))
}

// @Summary Create an app stack given app/env/stack
// @Tags    stacks
// @Accept  application/json
// @Param   payload body model.AppStack true "Specification of the stack"
// @Produce json
// @Success 200 {object} model.WrappedAppStack
// @Router  /v1/stacks/ [POST]
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
// @Router  /v1/stacks/ [DELETE]
// func (s *StackHandler) deleteAppStackHandler(ctx *fiber.Ctx) error {
// 	payload := getPayload[model.AppStackPayload](ctx)
// 	stack, err := s.stack.DeleteAppStack(payload)
// 	if err != nil {
// 		return response.ServerErrorResponse(ctx, err.Error())
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(wrapAppStack(stack))
// }

func wrapAppStacksWithCount(records []*model.AppStackResponse) model.WrappedAppStacksWithCount {
	return model.WrappedAppStacksWithCount{
		Records: records,
		Count:   len(records),
	}
}

// func wrapAppStack(record *model.AppStack) model.WrappedAppStack {
// 	return model.WrappedAppStack{
// 		Record: record,
// 	}
// }
