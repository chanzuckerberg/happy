package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/stretchr/testify/require"
)

type TestStackBackendECS struct{}

func (s *TestStackBackendECS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	result := model.MakeAppStack(payload.AppName, payload.Environment, "from-ecs")
	return []*model.AppStack{&result}, nil
}

type TestStackBackendEKS struct{}

func (s *TestStackBackendEKS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	result := model.MakeAppStack(payload.AppName, payload.Environment, "from-eks")
	return []*model.AppStack{&result}, nil
}

func TestGetFromBackendSuccess(t *testing.T) {
	testData := []struct {
		request           model.AppStackPayload
		expectedStackName string
	}{
		{
			request:           model.MakeAppStackPayload("testapp", "rdev", "", "czi-si", "us-west-2", "fargate", "", ""),
			expectedStackName: "from-ecs",
		},
		{
			request:           model.MakeAppStackPayload("testapp", "rdev", "", "czi-si", "us-west-2", "k8s", "testapp-rdev", "testapp-cluster"),
			expectedStackName: "from-eks",
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(r)
			err := db.AutoMigrate()
			r.NoError(err)

			stack := &Stack{
				ecs: &TestStackBackendECS{},
				eks: &TestStackBackendEKS{},
			}

			stacks, err := stack.GetAppStacks(context.Background(), tc.request)
			r.NoError(err)

			r.Equal(stacks[0].Stack, tc.expectedStackName)
		})
	}
}
