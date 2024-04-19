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
| [datadog_dashboard_json.stack_dashboard](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/dashboard_json) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_id"></a> [cluster\_id](#input\_cluster\_id) | Cluster ID | `string` | n/a | yes |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | Opsgenie Owner | `string` | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Stack name | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->