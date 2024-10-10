<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_datadog"></a> [datadog](#requirement\_datadog) | >= 3.20.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_datadog"></a> [datadog](#provider\_datadog) | >= 3.20.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [datadog_synthetics_test.test_api](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/synthetics_test) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage | `string` | n/a | yes |
| <a name="input_opsgenie_owner"></a> [opsgenie\_owner](#input\_opsgenie\_owner) | Opsgenie Owner | `string` | n/a | yes |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | Service name | `string` | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Stack name | `string` | n/a | yes |
| <a name="input_synthetic_url"></a> [synthetic\_url](#input\_synthetic\_url) | URL to run synthetic tests against | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | tags | `list(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END -->