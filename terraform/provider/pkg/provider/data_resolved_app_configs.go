package provider

import (
	"context"

	"github.com/chanzuckerberg/happy-shared/client"
	"github.com/chanzuckerberg/happy-shared/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResolvedAppConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: getResolvedAppConfigs,
		Schema: map[string]*schema.Schema{
			"appName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the app to fetch configs for",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The app environment to fetch configs for",
			},
			"stack": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the stack to fetch configs for",
			},

			"configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of config keys and values",
			},
		},
	}
}

func getResolvedAppConfigs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(client.HappyClient)

	appName := d.Get("app_name").(string)
	environment := d.Get("environment").(string)
	stack := d.Get("stack").(string)

	body := model.NewAppMetadata(appName, environment, stack)
	api.GetParsed("/v1/configs", body)

	return nil
}
