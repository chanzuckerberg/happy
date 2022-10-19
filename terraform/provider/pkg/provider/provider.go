package provider

import (
	"github.com/chanzuckerberg/happy-shared/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apiBaseUrl": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_BASE_URL", nil),
			},
			// TODO: add inputs needed for oidc
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"resolved_app_configs": ResolvedAppConfigs(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	apiBaseUrl := d.Get("apiBaseUrl").(string)
	api = client.NewHappyClient("happy-provider", "0.0.0", apiBaseUrl)

	return api, nil
}
