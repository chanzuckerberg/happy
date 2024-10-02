<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [kubernetes_manifest.linkerd_authorization_policy](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_manifest.linkerd_mesh_tls_authentication](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_manifest.linkerd_server](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allow_k6_operator"></a> [allow\_k6\_operator](#input\_allow\_k6\_operator) | A flag to allow the k6 operator to access this protected service | `bool` | `false` | no |
| <a name="input_allow_mesh_services"></a> [allow\_mesh\_services](#input\_allow\_mesh\_services) | A list of service/stack that we want to allow access to this protected service | <pre>list(object({<br>    service : optional(string, null),<br>    stack : optional(string, null),<br>    service_account_name : optional(string, null),<br>  }))</pre> | n/a | yes |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | n/a | yes |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service being protected | `string` | n/a | yes |
| <a name="input_labels"></a> [labels](#input\_labels) | Labels to apply to Linkerd CRDs | `map(string)` | n/a | yes |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | Name of the service being protected | `string` | n/a | yes |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | Port of the service being protected | `number` | n/a | yes |
| <a name="input_service_type"></a> [service\_type](#input\_service\_type) | Type of the service being protected | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->
