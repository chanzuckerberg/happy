<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_okta"></a> [okta](#provider\_okta) | ~> 4.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_happy_app"></a> [happy\_app](#module\_happy\_app) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth-head?okta-app-oauth-head-v0.1.0 | n/a |

## Resources

| Name | Type |
|------|------|
| [okta_app_group_assignments.happy_app](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/app_group_assignments) | resource |
| [okta_groups.teams](https://registry.terraform.io/providers/okta/okta/latest/docs/data-sources/groups) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The name of the happy application | `string` | n/a | yes |
| <a name="input_app_type"></a> [app\_type](#input\_app\_type) | The type of OAuth application. Valid values: `web`, `native`, `browser`, `service`. For SPA apps use `browser`. | `string` | `"web"` | no |
| <a name="input_aws_ssm_paths"></a> [aws\_ssm\_paths](#input\_aws\_ssm\_paths) | The name of the SSM paths for the client ID, secret, and other values produced from this app. | <pre>object({<br>    client_id     = string<br>    client_secret = string<br>    okta_idp_url  = string<br>    config_uri    = string<br>  })</pre> | <pre>{<br>  "client_id": "oauth2_proxy_client_id",<br>  "client_secret": "oauth2_proxy_client_secret",<br>  "config_uri": "oauth2_proxy_config_uri",<br>  "okta_idp_url": "oauth2_proxy_oidc_issuer_url"<br>}</pre> | no |
| <a name="input_base_domain"></a> [base\_domain](#input\_base\_domain) | The base domain to use for all the happy stacks. The default is app\_name.env.si.czi.technology | `string` | `"si.czi.technology"` | no |
| <a name="input_env"></a> [env](#input\_env) | The environment this happy application supports | `string` | n/a | yes |
| <a name="input_grant_types"></a> [grant\_types](#input\_grant\_types) | Additional grant types (authorization\_code is offered by default) | `list(string)` | <pre>[<br>  "authorization_code"<br>]</pre> | no |
| <a name="input_login_uri"></a> [login\_uri](#input\_login\_uri) | n/a | `string` | `""` | no |
| <a name="input_omit_secret"></a> [omit\_secret](#input\_omit\_secret) | Whether the provider should persist the application's secret to state. Your app's client\_secret will be recreated if this ever changes from true => false. | `bool` | `false` | no |
| <a name="input_redirect_uris"></a> [redirect\_uris](#input\_redirect\_uris) | n/a | `list(string)` | `[]` | no |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The component name this service is going to be deployed into | `string` | `"happy"` | no |
| <a name="input_teams"></a> [teams](#input\_teams) | The set of teams to give access to the Okta app | `set(string)` | n/a | yes |
| <a name="input_token_endpoint_auth_method"></a> [token\_endpoint\_auth\_method](#input\_token\_endpoint\_auth\_method) | Requested authentication method for the token endpoint. It can be set to `none`, `client_secret_post`, `client_secret_basic`, `client_secret_jwt`, `private_key_jwt`. To enable PKCE, set this to `none`. | `string` | `"client_secret_basic"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_oidc_config"></a> [oidc\_config](#output\_oidc\_config) | n/a |
<!-- END -->