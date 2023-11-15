locals {

  tasks = [for k, v in local.task_definitions : {
    "additionalNodeSelectors" = v.additional_node_selectors
    "additionalPodLabels"     = var.additional_pod_labels
    "awsIam" = {
      "roleArn" = v.aws_iam
    }
    "cmd" = v.cmd
    "env" = {
      "additionalEnvVars"               = merge(var.additional_env_vars, v.additional_env_vars, local.service_endpoints, local.stack_configs, local.context_env_vars)
      "additionalEnvVarsFromConfigMaps" = var.additional_env_vars_from_config_maps
      "additionalEnvVarsFromSecrets"    = var.additional_env_vars_from_secrets
    }
    "image" = {
      "platformArchitecture" = v.platform_architecture
      "pullPolicy"           = try(v.image_pull_policy, "IfNotPresent")
      "repository"           = v.image
      "tag"                  = lookup(var.image_tags, k, var.image_tag)
    }
    "name" = k
    "resources" = {
      "limits" = {
        "cpu"    = v.cpu
        "memory" = v.memory
      }
      "requests" = {
        "cpu"    = v.cpu
        "memory" = v.memory
      }
    }
    "schedule" = v.cron_schedule
    "suspend"  = v.is_cron_job ? false : true
    "volumes" = {
      "additionalVolumesFromConfigMaps" = [for v1 in var.additional_volumes_from_config_maps.items : {
        "mountPath" = "/var/${v1}"
        "name"      = v1
        "readOnly"  = true
      }]
      "additionalVolumesFromSecrets" = [for v1 in var.additional_volumes_from_secrets.items : {
        "mountPath" = "/var/${v1}"
        "name"      = v1
        "readOnly"  = true
      }]
    }
  }]

  services = [for k, v in local.service_definitions : {
    "additionalNodeSelectors" = v.additional_node_selectors
    "additionalPodLabels"     = var.additional_pod_labels
    "awsIam" = {
      "roleArn" = v.aws_iam
    }
    "certificateArn" = local.certificate_arn
    "args"           = v.args
    "cmd"            = v.cmd
    "env" = {
      "additionalEnvVars"               = merge(var.additional_env_vars, v.additional_env_vars, local.service_endpoints, local.stack_configs, local.context_env_vars)
      "additionalEnvVarsFromConfigMaps" = var.additional_env_vars_from_config_maps
      "additionalEnvVarsFromSecrets"    = var.additional_env_vars_from_secrets
    }
    "healthCheck" = {
      "initialDelaySeconds" = v.initial_delay_seconds
      "path"                = v.health_check_path
      "periodSeconds"       = v.period_seconds
    }
    "image" = {
      "repository"           = module.ecr[k].repository_url
      "tag"                  = lookup(var.image_tags, k, var.image_tag)
      "platformArchitecture" = v.platform_architecture
      "pullPolicy"           = try(v.image_pull_policy, "IfNotPresent")
      "scanOnPush"           = v.scan_on_push
      "tagMutability"        = v.tag_mutability
    }
    "name"             = k
    "regionalWafv2Arn" = local.regional_waf_arn
    "resources" = {
      "limits" = {
        "cpu"    = v.cpu
        "memory" = v.memory
      }
      "requests" = {
        "cpu"    = v.cpu_requests
        "memory" = v.memory_requests
      }
    }

    "routing" = {
      "alb" = {
        "loadBalancerAttributes" = [
          "idle_timeout.timeout_seconds=${v.alb_idle_timeout}",
        ]
        "targetGroup"    = v.group_name
        "targetGroupArn" = "" // TODO
        "securityGroups" = "" // TODO
      }
      "bypasses" = [
        (length(v.bypasses[k].methods) != 0 ? {
          field = "http-request-method"
          httpRequestMethodConfig = {
            Values = v.bypasses[k].methods
          }
        } : null),
        (length(v.bypasses[k].paths) != 0 ? {
          field = "path-pattern"
          pathPatternConfig = {
            Values = v.bypasses[k].paths
          }
        } : null)
      ]
      "groupName"    = v.group_name
      "hostMatch"    = v.host_match
      "method"       = var.routing_method
      "oidcConfig"   = local.oidc_config
      "path"         = v.path
      "port"         = v.port
      "priority"     = v.priority
      "scheme"       = v.scheme
      "serviceType"  = v.service_type
      "successCodes" = v.success_codes
    }
    "scaling" = {
      "cpuThresholdPercentage" = v.scaling_cpu_threshold_percentage
      "desiredCount"           = v.desired_count
      "maxCount"               = v.max_count
      "maxUnavailable"         = v.max_unavailable_count
    }
    "serviceMesh" = {
      "allowServices" = try(v.allow_mesh_services, [])
    }
    "sidecars" = [for k1, v1 in v.sidecars : {
      "healthCheck" = {
        "initialDelaySeconds" = v1.initial_delay_seconds
        "path"                = v1.health_check_path
        "periodSeconds"       = v1.period_seconds
      }
      "image" = {
        "repository" = v1.image
        "tag"        = v1.tag
      }
      "imagePullPolicy" = try(v1.image_pull_policy, "IfNotPresent")
      "name"            = k1
      "resources" = {
        "limits" = {
          "cpu"    = v1.cpu
          "memory" = v1.memory
        }
        "requests" = {
          "cpu"    = v1.cpu
          "memory" = v1.memory
        }
      }
      "routing" = {
        "port"   = v1.port
        "scheme" = v1.scheme
      }
    }]

    "volumes" = {
      "additionalVolumesFromConfigMaps" = [for v1 in var.additional_volumes_from_config_maps.items : {
        "mountPath" = "/var/${v1}"
        "name"      = v1
        "readOnly"  = true
      }]
      "additionalVolumesFromSecrets" = [for v1 in var.additional_volumes_from_secrets.items : {
        "mountPath" = "/var/${v1}"
        "name"      = v1
        "readOnly"  = true
      }]
    }
    "waitForSteadyState" = true
  }]

  values = {
    "stackName" = var.stack_name
    "aws" = {
      "cloudEnv" = {
        "databaseSubnetGroup" = local.cloud_env["database_subnet_group"]
        "databaseSubnets"     = local.cloud_env["database_subnets"]
        "privateSubnets"      = local.cloud_env["private_subnets"]
        "publicSubnets"       = local.cloud_env["public_subnets"]
        "vpcCidrBlock"        = local.cloud_env["vpc_cidr_block"]
        "vpcId"               = local.cloud_env["vpc_id"]
      }
      "dnsZone"   = local.secret["external_zone_name"]
      "region"    = data.aws_region.current.name
      "tags"      = local.tags
      "wafAclArn" = local.regional_waf_arn
    }
    "datadog" = {
      "enabled" = try(var.features["datadog"].enabled, false)
    }
    "deploymentStage" = var.deployment_stage
    "serviceMesh" = {
      "enabled" = try(var.features["service_mesh"].enabled, try(var.enable_service_mesh, false))
    }
    "services" = local.services
    "tasks"    = local.tasks
  }
}

resource "helm_release" "stack" {
  name       = var.app_name
  repository = "https://chanzuckerberg.github.io/happy-stack-helm/"
  chart      = "happy-stack"
  namespace  = var.k8s_namespace
  values     = [yamlencode(local.values)]
  wait       = true
}
