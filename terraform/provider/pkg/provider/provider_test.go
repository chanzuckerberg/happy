package provider

import (
	"context"
	"os"
	"testing"

	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/mock"
)

type APIMock struct {
	client.HappyConfigAPI
	mock.Mock
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func getTestProviders() (map[string]*schema.Provider, *APIMock) {
	happyProvider := Provider()
	apiMock := &APIMock{}
	happyProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		_, err := validateConfiguration(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		client := &APIClient{api: apiMock}
		return client, nil
	}
	providers := map[string]*schema.Provider{
		"happy": happyProvider,
	}
	return providers, apiMock
}

func testPreCheck(t *testing.T) {
	if err := os.Getenv("HAPPY_API_BASE_URL"); err == "" {
		t.Fatal("HAPPY_API_BASE_URL must be set for acceptance tests")
	}
}
