locals {
  tasks = [for k, v in var.tasks : {
    "additionalNodeSelectors" = v.additional_node_selectors
    "additionalPodLabels"     = var.additional_pod_labels
    "awsIam" = {
      "roleArn" = v.aws_iam
    }
    "cmd" = v.cmd
    "env" = {
      "additionalEnvVars"               = merge(var.additional_env_vars, v.additional_env_vars)
      "additionalEnvVarsFromConfigMaps" = var.additional_env_vars_from_config_maps
      "additionalEnvVarsFromSecrets"    = var.additional_env_vars_from_secrets
    }
    "image" = {
      "platformArchitecture" = v.platform_architecture
      "pullPolicy"           = try(v.image_pull_policy, "IfNotPresent")
      "repository"           = "blalbhal"
      "tag"                  = var.image_tag
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
      "additionalVolumesFromConfigMaps" = [for k1, v1 in var.additional_volumes_from_config_maps : {
        "mountPath" = v1.base_dir
        "name"      = k1
        "readOnly"  = true
      }]
      "additionalVolumesFromSecrets" = [for k1, v1 in var.additional_volumes_from_secrets : {
        "mountPath" = v1.base_dir
        "name"      = k1
        "readOnly"  = true
      }]
    }
  }]

  services = [for k, v in var.services : {
    "additionalNodeSelectors" = v.additional_node_selectors
    "additionalPodLabels"     = var.additional_pod_labels
    "awsIam" = {
      "roleArn" = v.aws_iam
    }
    "certificateArn" = "blahblahbs" // TODO
    "args"           = v.args
    "cmd"            = v.cmd
    "env" = {
      "additionalEnvVars"               = merge(var.additional_env_vars, v.additional_env_vars)
      "additionalEnvVarsFromConfigMaps" = var.additional_env_vars_from_config_maps
      "additionalEnvVarsFromSecrets"    = var.additional_env_vars_from_secrets
    }
    "healthCheck" = {
      "initialDelaySeconds" = v.initial_delay_seconds
      "path"                = v.health_check_path
      "periodSeconds"       = v.period_seconds
    }
    "image" = {
      "platformArchitecture" = v.platform_architecture
      "pullPolicy"           = try(v.image_pull_policy, "IfNotPresent")
      "repository"           = "blalbhal" // TODO
      "scanOnPush"           = v.scan_on_push
      "tag"                  = var.image_tag
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
          "idle_timeout.timeout_seconds=60", // TODO
        ]
        "securityGroup"  = "sg-123"                                                                 // TODO
        "targetGroup"    = "group1"                                                                 // TODO
        "targetGroupArn" = "arn:aws:elasticloadbalancing:us-west-2:00000000000:targetgroup/zzz/zzz" // TODO
      }
      "bypasses" = [ // TODO
        {
          "field" = "http-request-method" // TODO
          "httpRequestMethodConfig" = {
            "Values" = [ // TODO
              "GET",
              "OPTIONS",
            ]
          }
        },
        {
          "field" = "path-pattern" // TODO
          "pathPatternConfig" = {
            "Values" = [ // TODO
              "/blah",
              "/test/skip",
            ]
          }
        },
      ]
      "groupName" = "" // TODO
      "hostMatch" = "" // TODO
      "method"    = var.routing_method
      "oidcConfig" = {
        "authorizationEndpoint" = "" // TODO
        "issuer"                = "" // TODO
        "secretName"            = "" // TODO
        "tokenEndpoint"         = "" // TODO
        "userInfoEndpoint"      = "" // TODO
      }
      "path"         = v.path
      "port"         = v.port
      "priority"     = 4 // TODO
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
    "serviceEndpoints" = {} // TODO
    "serviceMesh" = {       // TODO
      "allowServices" = [   // TODO
        {
          "service"            = "service1" // TODO
          "serviceAccountName" = v.serviceAccountName
          "stack"              = "stack1" // TODO
        },
      ]
    }
    "sidecars" = [
      {
        "healthCheck" = {
          "initialDelaySeconds" = 30        // TODO
          "path"                = "/health" // TODO
          "periodSeconds"       = 3         // TODO
        }
        "image" = {                 // TODO
          "repository" = "blalbhal" // TODO
          "tag"        = "tag1"     // TODO
        }
        "imagePullPolicy"     = "IfNotPresent" // TODO
        "initialDelaySeconds" = 15             // TODO
        "name"                = "sidecar1"     // TODO
        "periodSeconds"       = 5              // TODO
        "resources" = {
          "limits" = {         // TODO
            "cpu"    = "100m"  // TODO
            "memory" = "100Mi" // TODO
          }
          "requests" = {
            "cpu"    = "10m"  // TODO
            "memory" = "10Mi" // TODO
          }
        }
        "routing" = {
          "port"   = 8080   // TODO
          "scheme" = "HTTP" // TODO
        }
      },
    ]
    "skipConfigInjection" = false // TODO
    "stackPrefix"         = ""    // TODO
    "volumes" = {
      "additionalVolumesFromConfigMaps" = [for k1, v1 in var.additional_volumes_from_config_maps : {
        "mountPath" = v1.base_dir
        "name"      = k1
        "readOnly"  = true
      }]
      "additionalVolumesFromSecrets" = [for k1, v1 in var.additional_volumes_from_secrets : {
        "mountPath" = v1.base_dir
        "name"      = k1
        "readOnly"  = true
      }]
    }
    "waitForSteadyState" = true
  }]

  values = {
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
      "enabled" = try(var.features["service_mesh"].enabled, false)
    }
    "services"  = local.services
    "stackName" = var.stack_name
    "tasks"     = local.tasks
  }
}

data "kubernetes_secret" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = var.k8s_namespace
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
