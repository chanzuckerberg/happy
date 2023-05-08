package aws

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

const testK8sFilePath = "../../config/testdata/test_k8s_config.yaml"
const testKubeconfig = "../../config/testdata/kubeconfig"

func TestK8SComputeBackend(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testK8sFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	r.Equal("eks", happyConfig.K8SConfig().AuthMethod)
	r.NotEmpty(happyConfig.K8SConfig().ClusterID)
	r.NotEmpty(happyConfig.K8SConfig().Namespace)

	secretsApi := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secretsApi.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

	stsPresignApi := interfaces.NewMockSTSPresignAPI(ctrl)
	stsPresignApi.EXPECT().PresignGetCallerIdentity(gomock.Any(), gomock.Any(), gomock.Any()).Return(&v4.PresignedHTTPRequest{URL: "", Method: "POST", SignedHeader: http.Header{}}, nil).AnyTimes()

	eksApi := interfaces.NewMockEKSAPI(ctrl)
	eksApi.EXPECT().DescribeCluster(gomock.Any(), gomock.Any()).Return(&eks.DescribeClusterOutput{Cluster: &ekstypes.Cluster{
		Arn:     aws.String("arn:aws:eks:us-west-2:1234567890:cluster/eks-cluster"),
		Name:    aws.String("eks-cluster"),
		Version: aws.String("1.23"),
		CertificateAuthority: &ekstypes.Certificate{
			Data: aws.String("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJZVENDQVFlZ0F3SUJBZ0lCS2pBS0JnZ3Foa2pPUFFRREFqQXBNUlF3RWdZRFZRUUxFd3RsYm1kcGJtVmwKY21sdVp6RVJNQThHQTFVRUF4TUlZMkZ6YUM1aGNIQXdIaGNOTnpBd01UQXhNREF3TURBMVdoY05OekF3TVRBeApNREF3TURFd1dqQXBNUlF3RWdZRFZRUUxFd3RsYm1kcGJtVmxjbWx1WnpFUk1BOEdBMVVFQXhNSVkyRnphQzVoCmNIQXdXVEFUQmdjcWhrak9QUUlCQmdncWhrak9QUU1CQndOQ0FBU2RhOENoa1FYeEdFTG5yVi9vQm5JQXgzZEQKb2NVT0pmZHo0cE9KVFA2ZFZRQjlVM1VCaVc1dVNYL01vT0QwTEw1ekczYlZ5TDNZNnBEd0t1WXZmTE5ob3lBdwpIakFjQmdOVkhSRUJBZjhFRWpBUWh3UUJBUUVCZ2doallYTm9MbUZ3Y0RBS0JnZ3Foa2pPUFFRREFnTklBREJGCkFpQXlISGcxTjZZRERRaVk5MjArY25JNVhTWndFR2hBdGI5UFlXTzhiTG1rY1FJaEFJMkNmRVpmM1Yvb2JtZFQKeXlhb0V1ZkxLVlhoclRRaFJmb2RUZWlnaTRSWAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t"),
		},
		Endpoint: aws.String("https://AABBCCDDEEFF.gr1.us-west-2.eks.amazonaws.com"),
	}}, nil).AnyTimes()

	// Can be used instead of the fake client if there's a need
	// kubernetesClient := interfaces.NewMockKubernetesAPI(ctrl)
	// corev1 := interfaces.NewMockKubernetesCoreV1API(ctrl)
	// kubernetesClient.EXPECT().CoreV1().Return(corev1).AnyTimes()
	// secretsGetterApi := interfaces.NewMockKubernetesSecretAPI(ctrl)
	// corev1.EXPECT().Secrets(gomock.Any()).Return(secretsGetterApi).AnyTimes()
	// secretsGetterApi.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&v1.Secret{
	// 	Data: map[string][]byte{
	// 		"integration_secret": []byte("{}"),
	// 	},
	// }, nil).AnyTimes()

	integrationSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "integration-secret",
			Namespace: happyConfig.K8SConfig().Namespace,
		},
		Immutable: new(bool),
		Data: map[string][]byte{
			"integration_secret": []byte(testVal),
		},
		Type: "",
	}

	cli := fake.NewSimpleClientset(integrationSecret)

	b, err := NewAWSBackend(ctx, happyConfig.GetEnvironmentContext(),
		WithAWSAccountID("1234567890"),
		WithSTSClient(stsApi),
		WithSTSPresignClient(stsPresignApi),
		WithEKSClient(eksApi),
		WithK8SClientCreator(func(config *rest.Config) (kubernetes.Interface, error) {
			return cli, nil
		}),
	)
	r.NoError(err)

	r.IsType(&K8SComputeBackend{}, b.GetComputeBackend())

	secret, secretArn, err := b.GetComputeBackend().GetIntegrationSecret(ctx)
	r.NoError(err)

	r.NotNil(secret)
	r.NotNil(secretArn)
	r.Empty(*secretArn)

	config, err := kube.CreateEKSConfig(ctx, eksApi, "eks-cluster")
	r.NoError(err)
	r.NotNil(config)

	token := kube.GetAuthToken(ctx, stsPresignApi, "eks-cluster")
	r.NotEmpty(token)
}

func TestK8SComputeBackendKubeconfig(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testK8sFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "stage",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	r.Equal("kubeconfig", happyConfig.K8SConfig().AuthMethod)
	r.NotEmpty(happyConfig.K8SConfig().ClusterID)
	r.NotEmpty(happyConfig.K8SConfig().Namespace)
	happyConfig.K8SConfig().KubeConfigPath = testKubeconfig

	secretsApi := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secretsApi.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

	stsPresignApi := interfaces.NewMockSTSPresignAPI(ctrl)
	stsPresignApi.EXPECT().PresignGetCallerIdentity(gomock.Any(), gomock.Any(), gomock.Any()).Return(&v4.PresignedHTTPRequest{URL: "", Method: "POST", SignedHeader: http.Header{}}, nil).AnyTimes()

	integrationSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "integration-secret",
			Namespace: happyConfig.K8SConfig().Namespace,
		},
		Immutable: new(bool),
		Data: map[string][]byte{
			"integration_secret": []byte(testVal),
		},
		Type: "",
	}

	cli := fake.NewSimpleClientset(integrationSecret)

	b, err := NewAWSBackend(ctx, happyConfig.GetEnvironmentContext(),
		WithAWSAccountID("1234567890"),
		WithSTSClient(stsApi),
		WithSTSPresignClient(stsPresignApi),
		WithK8SClientCreator(func(config *rest.Config) (kubernetes.Interface, error) {
			return cli, nil
		}),
	)
	r.NoError(err)

	r.IsType(&K8SComputeBackend{}, b.GetComputeBackend())

	secret, secretArn, err := b.GetComputeBackend().GetIntegrationSecret(ctx)
	r.NoError(err)

	r.NotNil(secret)
	r.NotNil(secretArn)
	r.Empty(*secretArn)

	config, err := kube.CreateK8SConfig(testKubeconfig, "cluster")
	r.NoError(err)
	r.NotNil(config)
}
