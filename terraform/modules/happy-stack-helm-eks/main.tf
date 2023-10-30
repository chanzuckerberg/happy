data "aws_region" "current" {}

locals {
  secret    = jsondecode(nonsensitive(data.kubernetes_secret.integration_secret.data.integration_secret))
  tags      = local.secret["tags"]
  cloud_env = local.secret["cloud_env"]

  waf_config       = lookup(local.secret, "waf_config", {})
  regional_waf_arn = lookup(local.waf_config, "arn", null)

  tasks = [for v in var.tasks : {
    "additionalNodeSelectors" = {}
    "additionalPodLabels"     = {}
    "awsIam" = {
      "roleArn" = "arn:aws:iam::00000000000:role/zzz/zzz"
    }
    "cmd" = [
      "./manage.py",
      "migrate",
    ]
    "env" = {
      "additionalEnvVars"               = []
      "additionalEnvVarsFromConfigMaps" = []
      "additionalEnvVarsFromSecrets"    = []
    }
    "image" = {
      "platformArchitecture" = "amd64"
      "pullPolicy"           = "IfNotPresent"
      "repository"           = "blalbhal"
      "tag"                  = "tag1"
    }
    "name" = "migrate"
    "resources" = {
      "limits" = {
        "cpu"    = "100m"
        "memory" = "100Mi"
      }
      "requests" = {
        "cpu"    = "10m"
        "memory" = "10Mi"
      }
    }
    "schedule" = "0 0 1 1 *"
    "suspend"  = true
    "volumes" = {
      "additionalVolumesFromConfigMaps" = [
        {
          "mountPath" = "blah"
          "name"      = "blah"
          "readOnly"  = true
        },
      ]
      "additionalVolumesFromSecrets" = [
        {
          "mountPath" = "blah2"
          "name"      = "blah2"
          "readOnly"  = true
        },
      ]
      "configMap" = {
        "items" = [
          {
            "key"  = "log_level"
            "path" = "log_level"
          },
        ]
        "name" = "log-config"
      }
    }
  }]

  services = [for v in var.services : {
    "additionalNodeSelectors" = {}
    "additionalPodLabels"     = {}
    "args"                    = []
    "awsIam" = {
      "roleArn" = "arn:aws:iam::00000000000:role/zzz/zzz"
    }
    "certificateArn" = "blahblahbs"
    "cmd"            = []
    "datadog" = {
      "createDashboard" = false
    }
    "env" = {
      "additionalEnvVars"               = []
      "additionalEnvVarsFromConfigMaps" = []
      "additionalEnvVarsFromSecrets"    = []
    }
    "healthCheck" = {
      "initialDelaySeconds" = 30
      "path"                = "/"
      "periodSeconds"       = 3
    }
    "image" = {
      "platformArchitecture" = "amd64"
      "pullPolicy"           = "IfNotPresent"
      "repository"           = "blalbhal"
      "scanOnPush"           = false
      "tag"                  = "tag1"
      "tagMutability"        = true
    }
    "name"             = "service2"
    "regionalWafv2Arn" = null
    "resources" = {
      "limits" = {
        "cpu"    = "100m"
        "memory" = "100Mi"
      }
      "requests" = {
        "cpu"    = "10m"
        "memory" = "10Mi"
      }
    }
    "routing" = {
      "alb" = {
        "loadBalancerAttributes" = [
          "idle_timeout.timeout_seconds=60",
        ]
        "securityGroup"  = "sg-123"
        "targetGroup"    = "group1"
        "targetGroupArn" = "arn:aws:elasticloadbalancing:us-west-2:00000000000:targetgroup/zzz/zzz"
      }
      "bypasses" = [
        {
          "field" = "http-request-method"
          "httpRequestMethodConfig" = {
            "Values" = [
              "GET",
              "OPTIONS",
            ]
          }
        },
        {
          "field" = "path-pattern"
          "pathPatternConfig" = {
            "Values" = [
              "/blah",
              "/test/skip",
            ]
          }
        },
      ]
      "groupName" = ""
      "hostMatch" = ""
      "method"    = "DOMAIN"
      "oidcConfig" = {
        "authorizationEndpoint" = ""
        "issuer"                = ""
        "secretName"            = ""
        "tokenEndpoint"         = ""
        "userInfoEndpoint"      = ""
      }
      "path"         = "/*"
      "port"         = 3000
      "priority"     = 4
      "scheme"       = "HTTP"
      "serviceName"  = ""
      "serviceType"  = "EXTERNAL"
      "successCodes" = "200-499"
    }
    "scaling" = {
      "cpuThresholdPercentage" = 80
      "desiredCount"           = 2
      "maxCount"               = 2
    }
    "serviceEndpoints" = {}
    "serviceMesh" = {
      "allowServices" = [
        {
          "service"            = "service1"
          "serviceAccountName" = "sa1"
          "stack"              = "stack1"
        },
      ]
    }
    "sidecars" = [
      {
        "healthCheck" = {
          "initialDelaySeconds" = 30
          "path"                = "/health"
          "periodSeconds"       = 3
        }
        "image" = {
          "repository" = "blalbhal"
          "tag"        = "tag1"
        }
        "imagePullPolicy"     = "IfNotPresent"
        "initialDelaySeconds" = 15
        "name"                = "sidecar1"
        "periodSeconds"       = 5
        "resources" = {
          "limits" = {
            "cpu"    = "100m"
            "memory" = "100Mi"
          }
          "requests" = {
            "cpu"    = "10m"
            "memory" = "10Mi"
          }
        }
        "routing" = {
          "port"   = 8080
          "scheme" = "HTTP"
        }
      },
    ]
    "skipConfigInjection" = false
    "stackPrefix"         = ""
    "volumes" = {
      "additionalVolumesFromConfigMaps" = [
        {
          "mountPath" = "blah"
          "name"      = "blah"
          "readOnly"  = true
        },
      ]
      "additionalVolumesFromSecrets" = [
        {
          "mountPath" = "blah2"
          "name"      = "blah2"
          "readOnly"  = true
        },
      ]
      "configMap" = {
        "items" = [
          {
            "key"  = "log_level"
            "path" = "log_level"
          },
        ]
        "name" = "log-config"
      }
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
  values     = local.values
  wait       = true
}
