<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_jwks"></a> [jwks](#requirement\_jwks) | 0.0.3 |
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 3.10 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_jwks"></a> [jwks](#provider\_jwks) | 0.0.3 |
| <a name="provider_okta"></a> [okta](#provider\_okta) | ~> 3.10 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_happy_app"></a> [happy\_app](#module\_happy\_app) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth | heathj/jwks |

## Resources

| Name | Type |
|------|------|
| [aws_kms_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key) | resource |
| [okta_app_group_assignments.happy_app](https://registry.terraform.io/providers/chanzuckerberg/okta/latest/docs/resources/app_group_assignments) | resource |
| [aws_kms_public_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_public_key) | data source |
| [jwks_from_key.jwks](https://registry.terraform.io/providers/iwarapter/jwks/0.0.3/docs/data-sources/from_key) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The name of the happy application | `string` | n/a | yes |
| <a name="input_aws_ssm_paths"></a> [aws\_ssm\_paths](#input\_aws\_ssm\_paths) | The name of the SSM paths for the client ID, secret, and other values produced from this app. | <pre>object({<br>    client_id     = string<br>    client_secret = string<br>    okta_idp_url  = string<br>    config_uri    = string<br>  })</pre> | <pre>{<br>  "client_id": "oauth2_proxy_client_id",<br>  "client_secret": "oauth2_proxy_client_secret",<br>  "config_uri": "oauth2_proxy_config_uri",<br>  "okta_idp_url": "oauth2_proxy_oidc_issuer_url"<br>}</pre> | no |
| <a name="input_login_uri"></a> [login\_uri](#input\_login\_uri) | n/a | `string` | `""` | no |
| <a name="input_rbac_role_mapping"></a> [rbac\_role\_mapping](#input\_rbac\_role\_mapping) | n/a | `map(list(string))` | `{}` | no |
| <a name="input_redirect_uris"></a> [redirect\_uris](#input\_redirect\_uris) | n/a | `list(string)` | `[]` | no |
| <a name="input_scope_name"></a> [scope\_name](#input\_scope\_name) | The name of the custom scope that allows the service account to authenticate with Client Credentials flow. | `string` | `"service_account"` | no |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The component name this service is going to be deployed into | `string` | `"happy"` | no |
| <a name="input_teams"></a> [teams](#input\_teams) | The set of teams to give access to the Okta app | `set(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END -->