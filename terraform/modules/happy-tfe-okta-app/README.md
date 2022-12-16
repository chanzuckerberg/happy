<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 3.10 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_okta"></a> [okta](#provider\_okta) | ~> 3.10 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_happy_apps"></a> [happy\_apps](#module\_happy\_apps) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth | c52bd57 |

## Resources

| Name | Type |
|------|------|
| [okta_app_group_assignments.happy_app](https://registry.terraform.io/providers/chanzuckerberg/okta/latest/docs/resources/app_group_assignments) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The name of the happy application | `string` | n/a | yes |
| <a name="input_app_type"></a> [app\_type](#input\_app\_type) | The type of OAuth application. Valid values: `web`, `native`, `browser`, `service`. For SPA apps use `browser`. | `string` | `"web"` | no |
| <a name="input_aws_ssm_paths"></a> [aws\_ssm\_paths](#input\_aws\_ssm\_paths) | The name of the SSM paths for the client ID, secret, and other values produced from this app. | <pre>object({<br>    client_id     = string<br>    client_secret = string<br>    okta_idp_url  = string<br>    config_uri    = string<br>  })</pre> | <pre>{<br>  "client_id": "oauth2_proxy_client_id",<br>  "client_secret": "oauth2_proxy_client_secret",<br>  "config_uri": "oauth2_proxy_config_uri",<br>  "okta_idp_url": "oauth2_proxy_oidc_issuer_url"<br>}</pre> | no |
| <a name="input_envs"></a> [envs](#input\_envs) | The environments this happy application supports | `set(string)` | n/a | yes |
| <a name="input_grant_types"></a> [grant\_types](#input\_grant\_types) | Additional grant types (authorization\_code is offered by default) | `list(string)` | <pre>[<br>  "authorization_code"<br>]</pre> | no |
| <a name="input_login_uri"></a> [login\_uri](#input\_login\_uri) | n/a | `string` | `""` | no |
| <a name="input_omit_secret"></a> [omit\_secret](#input\_omit\_secret) | n/a | `bool` | `false` | no |
| <a name="input_redirect_uris"></a> [redirect\_uris](#input\_redirect\_uris) | n/a | `list(string)` | `[]` | no |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The component name this service is going to be deployed into | `string` | `"happy"` | no |
| <a name="input_teams"></a> [teams](#input\_teams) | The set of teams to give access to the Okta app | `set(string)` | n/a | yes |
| <a name="input_token_endpoint_auth_method"></a> [token\_endpoint\_auth\_method](#input\_token\_endpoint\_auth\_method) | Requested authentication method for the token endpoint. It can be set to `none`, `client_secret_post`, `client_secret_basic`, `client_secret_jwt`, `private_key_jwt`. To enable PKCE, set this to `none`. | `string` | `"client_secret_basic"` | no |

## Outputs

No outputs.
<!-- END -->