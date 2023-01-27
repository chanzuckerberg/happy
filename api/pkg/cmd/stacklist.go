package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/api/pkg/backend"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/sirupsen/logrus"
)

type StacklistIface interface {
	GetAppStacks(model.AppStackPayload2) ([]*model.AppStack, error)
}

type Stacklist struct {
	k8s backend.StacklistK8S
	ecs backend.StacklistBackendECS
}

func MakeStacklist(ctx context.Context) *Stacklist {
	// k8s := backend.MakeStacklistBackendK8S(ctx)

	return &Stacklist{
		k8s: backend.StacklistK8S{},
	}
}
func (s *Stacklist) GetAppStacks(ctx context.Context, payload model.AppStackPayload2) ([]*model.AppStack, error) {
	switch payload.TaskLaunchType {
	case "k8s":
		return s.k8s.GetAppStacks(ctx, payload)
	case "fargate":
		return s.ecs.GetAppStacks(payload)
	default:
		logrus.Fatal("Configuration did not provide valid database driver and data_source_name")
	}
	return nil, nil
}
