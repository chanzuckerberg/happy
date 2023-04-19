package client

import "github.com/chanzuckerberg/happy/shared/model"

type HappyStackAPI interface {
	ListStacks(req model.AppStackPayload) (model.WrappedAppStacksWithCount, error)
}

func (c *HappyClient) ListStacks(req model.AppStackPayload) (model.WrappedAppStacksWithCount, error) {
	result := model.WrappedAppStacksWithCount{}
	err := c.GetParsed("/v1/stacks", req, &result)
	return result, err
}
