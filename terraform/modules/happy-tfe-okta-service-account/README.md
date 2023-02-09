<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_jwks"></a> [jwks](#requirement\_jwks) | 0.0.3 |
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 3.10 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_jwks"></a> [jwks](#provider\_jwks) | 0.0.3 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_service_user"></a> [service\_user](#module\_service\_user) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth | v0.233.0 |

## Resources

| Name | Type |
|------|------|
| [aws_kms_alias.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_alias) | resource |
| [aws_kms_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key) | resource |
| [aws_kms_public_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_public_key) | data source |
| [jwks_from_key.jwks](https://registry.terraform.io/providers/iwarapter/jwks/0.0.3/docs/data-sources/from_key) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The name of the happy application | `string` | n/a | yes |
| <a name="input_aws_ssm_paths"></a> [aws\_ssm\_paths](#input\_aws\_ssm\_paths) | The name of the SSM paths for the client ID, secret, and other values produced from this app. | <pre>object({<br>    client_id     = string<br>    client_secret = string<br>    okta_idp_url  = string<br>    config_uri    = string<br>  })</pre> | <pre>{<br>  "client_id": "oauth2_proxy_client_id",<br>  "client_secret": "oauth2_proxy_client_secret",<br>  "config_uri": "oauth2_proxy_config_uri",<br>  "okta_idp_url": "oauth2_proxy_oidc_issuer_url"<br>}</pre> | no |
| <a name="input_env"></a> [env](#input\_env) | The environment of this happy application | `string` | n/a | yes |
| <a name="input_rbac_role_mapping"></a> [rbac\_role\_mapping](#input\_rbac\_role\_mapping) | The roles that will be created as claims to access tokens for users authenticating to this application | `map(list(string))` | `{}` | no |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The component name this service is going to be deployed into | `string` | `"happy"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_kms_key_id"></a> [kms\_key\_id](#output\_kms\_key\_id) | n/a |
| <a name="output_oidc_config"></a> [oidc\_config](#output\_oidc\_config) | n/a |
<!-- END -->