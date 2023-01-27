package backend

import "github.com/chanzuckerberg/happy/shared/model"

type StacklistBackendECS struct{}

func (s *StacklistBackendECS) GetAppStacks(model.AppStackPayload2) ([]*model.AppStack, error) {
	return []*model.AppStack{}, nil
}
