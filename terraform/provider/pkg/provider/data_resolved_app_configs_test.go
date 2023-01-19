package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func (a *APIMock) ListConfigs(appName, env, stack string) (model.WrappedResolvedAppConfigsWithCount, error) {
	args := a.Called(appName, env, stack)
	output := args.Get(0).(model.WrappedResolvedAppConfigsWithCount)
	return output, args.Error(1)
}

func generateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func exportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func stashEnv() []string {
	env := os.Environ()
	os.Clearenv()
	return env
}

func popEnv(env []string) {
	os.Clearenv()

	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		k, v := p[0], ""
		if len(p) > 1 {
			v = p[1]
		}
		os.Setenv(k, v)
	}
}

func TestGetResolvedAppConfigsSucceed(t *testing.T) {
	// to make sure local environment doesn't mess with tests
	oldEnv := stashEnv()
	defer popEnv(oldEnv)

	r := require.New(t)
	providers, apiMock := getTestProviders()
	appName := "test-app"
	env := "rdev"
	stack := "foo"
	output := model.WrappedResolvedAppConfigsWithCount{
		Records: []*model.ResolvedAppConfig{
			{
				AppConfig: model.AppConfig{
					AppConfigPayload: *model.NewAppConfigPayload(appName, env, stack, "key1", "val1"),
				},
				Source: "stack",
			},
			{
				AppConfig: model.AppConfig{
					AppConfigPayload: *model.NewAppConfigPayload(appName, env, "", "key2", "val2"),
				},
				Source: "environment",
			},
		},
		Count: 1,
	}
	apiMock.On("ListConfigs", appName, env, stack).Return(output, nil)

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
				Config: testResolvedAppConfigsData(appName, env, stack),
				Check: func(s *terraform.State) error {
					configs := s.RootModule().Outputs["app_configs"].Value

					r.ElementsMatch(configs, []map[string]interface{}{
						{
							"key":    "key1",
							"value":  "val1",
							"source": "stack",
						},
						{
							"key":    "key2",
							"value":  "val2",
							"source": "environment",
						},
					})

					return nil
				},
			},
		},
	})
}

func testResolvedAppConfigsData(appName, env, stack string) string {
	return fmt.Sprintf(`
		data "happy_resolved_app_configs" "configs" {
			app_name    = "%s"
			environment = "%s"
			stack       = "%s"
		}

		output "app_configs" {
			value = data.happy_resolved_app_configs.configs.app_configs
		}
	`, appName, env, stack)
}
