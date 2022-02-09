package backend

// import (
// 	"testing"

// 	"github.com/aws/aws-sdk-go/service/sts"
// 	cziAWS "github.com/chanzuckerberg/go-misc/aws"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetUserName(t *testing.T) {
// 	r := require.New(t)
// 	ctrl := gomock.NewController(t)
// 	client := cziAWS.Client{}
// 	_, mock := client.WithMockSTS(ctrl)

// 	testVal := "test_user_name"
// 	mock.EXPECT().GetCallerIdentity(gomock.Any()).Return(&sts.GetCallerIdentityOutput{
// 		Arn: &testVal,
// 	}, nil)

// 	awsSecretMgr := GetAwsStsWithClient(mock)
// 	userName, err := awsSecretMgr.GetUserName()
// 	r.NoError(err)
// 	r.Equal(userName, "test_user_name")
// }
