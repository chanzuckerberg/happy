package provider

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResolvedAppConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: getResolvedAppConfigs,
		Schema: map[string]*schema.Schema{
			"app_name": {
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

			"app_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of config keys and values",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getResolvedAppConfigs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*APIClient)
	api := client.api

	appName := d.Get("app_name").(string)
	environment := d.Get("environment").(string)
	stack := d.Get("stack").(string)

	diags := diag.Diagnostics{}

	result, err := api.ListConfigs(appName, environment, stack)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(model.NewAppMetadata(appName, environment, stack).String())
	if err = d.Set("app_configs", getRecords(result)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func getRecords(result model.WrappedResolvedAppConfigsWithCount) []map[string]string {
	records := []map[string]string{}

	for _, config := range result.Records {
		records = append(records, map[string]string{
			"key":    config.Key,
			"value":  config.Value,
			"source": config.Source,
		})
	}

	return records
}
