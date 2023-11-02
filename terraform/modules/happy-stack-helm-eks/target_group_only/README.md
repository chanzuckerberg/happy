<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.23 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.23 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_lb_listener_rule.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_rule) | resource |
| [aws_lb_target_group.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group) | resource |
| [random_pet.this](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [aws_lb.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/lb) | data source |
| [aws_lb_listener.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/lb_listener) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | The path to use for the health check | `string` | `"/health"` | no |
| <a name="input_routing"></a> [routing](#input\_routing) | Routing configuration for the ingress | <pre>object({<br>    method : optional(string, "DOMAIN")<br>    host_match : string<br>    additional_hostnames : optional(set(string), [])<br>    group_name : string<br>    alb : optional(object({<br>      name : string,<br>      listener_port : number,<br>    }), null)<br>    priority : number<br>    path : optional(string, "/*")<br>    service_name : string<br>    port : number<br>    service_port : number<br>    alb_idle_timeout : optional(number, 60) // in seconds<br>    service_scheme : optional(string, "HTTP")<br>    scheme : optional(string, "HTTP")<br>    success_codes : optional(string, "200-499")<br>    service_type : string<br>    service_mesh : bool<br>    allow_mesh_services : optional(list(object({<br>      service : optional(string, null),<br>      stack : optional(string, null),<br>      service_account_name : optional(string, null),<br>    })), null)<br>    oidc_config : optional(object({<br>      issuer : string<br>      authorizationEndpoint : string<br>      tokenEndpoint : string<br>      userInfoEndpoint : string<br>      secretName : string<br>      }), {<br>      issuer                = ""<br>      authorizationEndpoint = ""<br>      tokenEndpoint         = ""<br>      userInfoEndpoint      = ""<br>      secretName            = ""<br>    })<br>    bypasses : optional(map(object({<br>      paths   = optional(set(string), [])<br>      methods = optional(set(string), [])<br>    })))<br>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_security_groups"></a> [security\_groups](#output\_security\_groups) | n/a |
<!-- END -->