<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_jwks"></a> [jwks](#requirement\_jwks) | 0.0.3 |
| <a name="requirement_okta"></a> [okta](#requirement\_okta) | ~> 4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |
| <a name="provider_jwks"></a> [jwks](#provider\_jwks) | 0.0.3 |
| <a name="provider_okta"></a> [okta](#provider\_okta) | ~> 4.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_params"></a> [params](#module\_params) | git@github.com:chanzuckerberg/cztack//aws-ssm-params-writer | v0.53.2 |

## Resources

| Name | Type |
|------|------|
| [aws_kms_alias.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_alias) | resource |
| [aws_kms_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key) | resource |
| [okta_app_oauth.app](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/app_oauth) | resource |
| [okta_auth_server.auth_server](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/auth_server) | resource |
| [okta_auth_server_policy.policy](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/auth_server_policy) | resource |
| [okta_auth_server_policy_rule.rule](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/auth_server_policy_rule) | resource |
| [okta_auth_server_scope.scope](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/auth_server_scope) | resource |
| [aws_kms_public_key.service_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_public_key) | data source |
| [jwks_from_key.jwks](https://registry.terraform.io/providers/iwarapter/jwks/0.0.3/docs/data-sources/from_key) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_okta_tenant"></a> [okta\_tenant](#input\_okta\_tenant) | The Okta tenant to create the authorization server. | `string` | `"czi-prod"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_app"></a> [app](#output\_app) | n/a |
| <a name="output_authz_server"></a> [authz\_server](#output\_authz\_server) | n/a |
| <a name="output_kms_key_id"></a> [kms\_key\_id](#output\_kms\_key\_id) | n/a |
| <a name="output_oidc_config"></a> [oidc\_config](#output\_oidc\_config) | n/a |
<!-- END -->