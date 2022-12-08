<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 3.10 |
| <a name="requirement_required_version"></a> [required\_version](#requirement\_required\_version) | >= 1.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_okta"></a> [okta](#provider\_okta) | ~> 3.10 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_happy_apps"></a> [happy\_apps](#module\_happy\_apps) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth | v0.202.0 |

## Resources

| Name | Type |
|------|------|
| [okta_app_group_assignments.happy_app](https://registry.terraform.io/providers/chanzuckerberg/okta/latest/docs/resources/app_group_assignments) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The name of the happy application | `string` | n/a | yes |
| <a name="input_aws_ssm_paths"></a> [aws\_ssm\_paths](#input\_aws\_ssm\_paths) | The name of the SSM paths for the client ID, secret, and other values produced from this app. | <pre>object({<br>    client_id     = string<br>    client_secret = string<br>    okta_idp_url  = string<br>    config_uri    = string<br>  })</pre> | <pre>{<br>  "client_id": "oauth2_proxy_client_id",<br>  "client_secret": "oauth2_proxy_client_secret",<br>  "config_uri": "oauth2_proxy_config_uri",<br>  "okta_idp_url": "oauth2_proxy_oidc_issuer_url"<br>}</pre> | no |
| <a name="input_envs"></a> [envs](#input\_envs) | The environments this happy application supports | `set(string)` | n/a | yes |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The component name this service is going to be deployed into | `string` | `"happy"` | no |
| <a name="input_teams"></a> [teams](#input\_teams) | The set of teams to give access to the Okta app | `set(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END -->