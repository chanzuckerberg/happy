package provider

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Stacklist() *schema.Resource {
	return &schema.Resource{
		ReadContext: getStacklist,
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
			"aws_profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS profile where the app/env is deployed",
			},
			"aws_region": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The AWS region where the app/env is deployed",
				Default:     "us-west-2",
			},
			"task_launch_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "'fargate' or 'k8s', depending where the app/env is deployed",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					launchType := v.(string)
					if launchType != "fargate" && launchType != "k8s" {
						return diag.Diagnostics{
							diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Invalid value",
								Detail:   "Must be either 'fargate' or 'k8s'",
							},
						}
					}

					return diag.Diagnostics{}
				},
			},
			"k8s_namespace": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Description:  "The k8s namespace where the app/env is deployed",
				RequiredWith: []string{"k8s_namespace", "k8s_cluster_id"},
			},
			"k8s_cluster_id": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Description:  "The k8s cluster ID where the app/env is deployed",
				RequiredWith: []string{"k8s_namespace", "k8s_cluster_id"},
			},

			"stacklist": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of stack names for the given env",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getStacklist(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*APIClient)
	api := client.api

	appName := d.Get("app_name").(string)
	environment := d.Get("environment").(string)
	stack := ""
	awsProfile := d.Get("aws_profile").(string)
	awsRegion := d.Get("aws_region").(string)
	launchType := d.Get("task_launch_type").(string)
	k8sNamespace := d.Get("k8s_namespace").(string)
	k8sClusterId := d.Get("k8s_cluster_id").(string)

	if launchType == "k8s" && (k8sNamespace == "" || k8sClusterId == "") {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing k8s values",
				Detail:   "'k8s_namespace' and 'k8s_cluster_id' must be provided when 'task_launch_type' is 'k8s'",
			},
		}
	}

	request := model.MakeAppStackPayload(appName, environment, stack, awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId)
	result, err := api.ListStacks(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(model.NewAppMetadata(appName, environment, stack).String())
	if err = d.Set("stacklist", getStackNames(result)); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func getStackNames(result model.WrappedAppStacksWithCount) []string {
	stacklist := []string{}
	for _, record := range result.Records {
		stacklist = append(stacklist, record.Stack)
	}
	return stacklist
}
