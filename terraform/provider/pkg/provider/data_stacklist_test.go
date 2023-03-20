package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func (a *APIMock) ListStacks(request model.AppStackPayload) (model.WrappedAppStacksWithCount, error) {
	args := a.Called(request)
	output := args.Get(0).(model.WrappedAppStacksWithCount)
	return output, args.Error(1)
}

func TestGetFargateStacklistSucceed(t *testing.T) {
	// to make sure local environment doesn't mess with tests
	oldEnv := stashEnv()
	defer popEnv(oldEnv)

	r := require.New(t)
	providers, apiMock := getTestProviders()
	appName := "testapp"
	env := "rdev"
	awsProfile := "czi-playground"
	awsRegion := "us-west-2"
	launchType := "fargate"
	k8sNamespace := ""
	k8sClusterId := ""

	stacks := []string{"foo", "bar", "baz"}
	records := []*model.AppStack{}
	for _, s := range stacks {
		records = append(records, &model.AppStack{
			AppMetadata: *model.NewAppMetadata(appName, env, s),
		})
	}

	output := model.WrappedAppStacksWithCount{
		Records: records,
		Count:   1,
	}
	apiMock.On("ListStacks", model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId)).Return(output, nil)

	private, _ := generateRsaKeyPair()
	pemString := exportRsaPrivateKeyAsPemStr(private)
	os.Setenv("TF_ACC", "yes")
	os.Setenv("HAPPY_API_BASE_URL", "https://fake.happy-api.io")
	os.Setenv("HAPPY_API_PRIVATE_KEY", pemString)
	os.Setenv("HAPPY_API_OIDC_ISSUER", "fake-issuer")
	os.Setenv("HAPPY_API_OIDC_AUTHZ_ID", "fake-authz-id")
	os.Setenv("HAPPY_API_OIDC_SCOPE", "fake-scope")
	os.Setenv("HAPPY_API_ASSUME_ROLE_ARN", "fake-role")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testPreCheck(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testFargateStacklistData(appName, env, awsProfile, awsRegion, launchType),
				Check: func(s *terraform.State) error {
					stacks := s.RootModule().Outputs["stacks"].Value
					r.ElementsMatch(stacks, stacks)
					return nil
				},
			},
		},
	})
}

func TestGetK8sStacklistSucceed(t *testing.T) {
	// to make sure local environment doesn't mess with tests
	oldEnv := stashEnv()
	defer popEnv(oldEnv)

	r := require.New(t)
	providers, apiMock := getTestProviders()
	appName := "testapp"
	env := "rdev"
	awsProfile := "czi-playground"
	awsRegion := "us-west-2"
	launchType := "k8s"
	k8sNamespace := "test-ns"
	k8sClusterId := "test-cluster"

	stacks := []string{"foo", "bar", "baz"}
	records := []*model.AppStack{}
	for _, s := range stacks {
		records = append(records, &model.AppStack{
			AppMetadata: *model.NewAppMetadata(appName, env, s),
		})
	}

	output := model.WrappedAppStacksWithCount{
		Records: records,
		Count:   1,
	}
	apiMock.On("ListStacks", model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId)).Return(output, nil)

	private, _ := generateRsaKeyPair()
	pemString := exportRsaPrivateKeyAsPemStr(private)
	os.Setenv("TF_ACC", "yes")
	os.Setenv("HAPPY_API_BASE_URL", "https://fake.happy-api.io")
	os.Setenv("HAPPY_API_PRIVATE_KEY", pemString)
	os.Setenv("HAPPY_API_OIDC_ISSUER", "fake-issuer")
	os.Setenv("HAPPY_API_OIDC_AUTHZ_ID", "fake-authz-id")
	os.Setenv("HAPPY_API_OIDC_SCOPE", "fake-scope")
	os.Setenv("HAPPY_API_ASSUME_ROLE_ARN", "fake-role")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testPreCheck(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testK8sStacklistData(appName, env, awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId),
				Check: func(s *terraform.State) error {
					stacks := s.RootModule().Outputs["stacks"].Value
					r.ElementsMatch(stacks, stacks)
					return nil
				},
			},
		},
	})
}

func TestGetStacklistArgumentErrors(t *testing.T) {
	providers, _ := getTestProviders()
	appName := "testapp"
	env := "rdev"
	awsProfile := "czi-playground"
	awsRegion := "us-west-2"

	testData := []struct {
		payload      model.AppStackPayload
		errorMessage string
	}{
		{
			payload:      model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, "foogate", "", ""),
			errorMessage: "Must be either 'fargate' or 'k8s'",
		},
		{
			payload:      model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, "k8s", "", ""),
			errorMessage: "'k8s_namespace' and 'k8s_cluster_id' must be provided when 'task_launch_type' is 'k8s'",
		},
		{
			payload:      model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, "k8s", "foo", ""),
			errorMessage: "'k8s_namespace' and 'k8s_cluster_id' must be provided when 'task_launch_type' is 'k8s'",
		},
		{
			payload:      model.MakeAppStackPayload(appName, env, "", awsProfile, awsRegion, "k8s", "", "foo"),
			errorMessage: "'k8s_namespace' and 'k8s_cluster_id' must be provided when 'task_launch_type' is 'k8s'",
		},
	}

	// to make sure local environment doesn't mess with tests
	oldEnv := stashEnv()
	defer popEnv(oldEnv)

	for idx, testCase := range testData {
		tc := testCase

		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			private, _ := generateRsaKeyPair()
			pemString := exportRsaPrivateKeyAsPemStr(private)
			os.Setenv("TF_ACC", "yes")
			os.Setenv("HAPPY_API_BASE_URL", "https://fake.happy-api.io")
			os.Setenv("HAPPY_API_PRIVATE_KEY", pemString)
			os.Setenv("HAPPY_API_OIDC_ISSUER", "fake-issuer")
			os.Setenv("HAPPY_API_OIDC_AUTHZ_ID", "fake-authz-id")
			os.Setenv("HAPPY_API_OIDC_SCOPE", "fake-scope")
			os.Setenv("HAPPY_API_ASSUME_ROLE_ARN", "fake-role")

			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testPreCheck(t) },
				Providers: providers,
				Steps: []resource.TestStep{
					{
						PlanOnly: true,
						Config: testFargateStacklistData(
							tc.payload.AppName,
							tc.payload.Environment,
							tc.payload.AwsProfile,
							tc.payload.AwsRegion,
							tc.payload.TaskLaunchType,
						),
						// replace spaces so regex will match arbitrary line breaks
						ExpectError: regexp.MustCompile(strings.Replace(tc.errorMessage, " ", "\\s", -1)),
					},
				},
			})
		})
	}
}

func testFargateStacklistData(appName, env, awsProfile, awsRegion, launchType string) string {
	return fmt.Sprintf(`
		data "happy_stacklist" "stacks" {
			app_name         = "%s"
			environment      = "%s"
			aws_profile      = "%s"
			task_launch_type = "%s"
		}

		output "stacks" {
			value = data.happy_stacklist.stacks.stacklist
		}
	`, appName, env, awsProfile, launchType)
}

func testK8sStacklistData(appName, env, awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId string) string {
	return fmt.Sprintf(`
		data "happy_stacklist" "stacks" {
			app_name         = "%s"
			environment      = "%s"
			aws_profile      = "%s"
			task_launch_type = "%s"
			k8s_namespace    = "%s"
			k8s_cluster_id   = "%s"
		}

		output "stacks" {
			value = data.happy_stacklist.stacks.stacklist
		}
	`, appName, env, awsProfile, launchType, k8sNamespace, k8sClusterId)
}
