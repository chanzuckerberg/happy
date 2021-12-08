package backend

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetParameter(t *testing.T) {
	r := require.New(t)
	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSSM(ctrl)
	testVal := "test_param_val"
	mock.EXPECT().GetParameter(gomock.Any()).Return(&ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Value: &testVal,
		},
	},
		nil)

	awsBackend := GetAwsBackendWithClient(mock)
	out, err := awsBackend.GetParameter("test_param_path")
	r.Nil(err)
	r.Equal("test_param_val", *out)
}

func TestAddParams(t *testing.T) {
	r := require.New(t)
	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSSM(ctrl)
	mock.EXPECT().PutParameter(gomock.Any()).Return(&ssm.PutParameterOutput{}, nil)

	awsBackend := GetAwsBackendWithClient(mock)
	err := awsBackend.AddParams("test_param_name", "test_param_val")
	r.Nil(err)
}
