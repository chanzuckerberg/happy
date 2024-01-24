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
| [kubernetes_ingress_v1.ingress](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_host_match"></a> [host\_match](#input\_host\_match) | Host header value to match when routing to the service | `string` | n/a | yes |
| <a name="input_host_path"></a> [host\_path](#input\_host\_path) | Path value to match when routing to the service | `string` | `"/"` | no |
| <a name="input_ingress_name"></a> [ingress\_name](#input\_ingress\_name) | Name of the ingress resource | `string` | n/a | yes |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_labels"></a> [labels](#input\_labels) | Labels to apply to ingress resource | `map(string)` | n/a | yes |
| <a name="input_sticky_sessions"></a> [sticky\_sessions](#input\_sticky\_sessions) | Sticky session configuration | <pre>object({<br>    enabled          = optional(bool, true),<br>    duration_seconds = optional(number, 600),<br>    cookie_name      = optional(string, "happy_sticky_session"),<br>  })</pre> | `{}` | no |
| <a name="input_target_service_name"></a> [target\_service\_name](#input\_target\_service\_name) | Name of destination service that the ingress should route to | `string` | n/a | yes |
| <a name="input_target_service_port"></a> [target\_service\_port](#input\_target\_service\_port) | Port of destination service that the ingress should route to | `string` | n/a | yes |
| <a name="input_timeout"></a> [timeout](#input\_timeout) | Timeout for the ingress resource | `number` | `60` | no |

## Outputs

No outputs.
<!-- END -->