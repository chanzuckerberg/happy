package provider

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/terraform/provider/pkg/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIClient struct {
	api client.HappyConfigAPI
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_BASE_URL", nil),
			},
			// TODO: add inputs needed for oidc
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"happy_resolved_app_configs": ResolvedAppConfigs(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiBaseUrl := d.Get("api_base_url").(string)
	api := client.NewHappyClient("happy-provider", version.ProviderVersion, apiBaseUrl)
	return &APIClient{api: api}, nil
}
