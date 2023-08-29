---
parent: Config
layout: default
has_toc: true
---

# Config.json
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

# Config Files

### .happy/config.json

This file is the primary configuration file for your happy application. It contains information about the various environments and AWS
accounts you want to deploy, the services you want to deploy, as well as deployment options such as memory, CPU and health check endpoints
of your services. Here is an example config.json file:

~~~json
{
    "default_env": "rdev",
    "app": "scruffy",
    "environments": {
        "prod": {
            "aws_profile": "czi-si",
            "k8s": {
                "namespace": "scruffy-prod-happy-happy-env",
                "cluster_id": "hapi-prod-eks",
                "auth_method": "eks"
            },
            "terraform_directory": ".happy/terraform/envs/prod",
            "task_launch_type": "k8s"
        },
        "rdev": {
            "aws_profile": "czi-si",
            "k8s": {
                "namespace": "scruffy-rdev-happy-happy-env",
                "cluster_id": "hapi-rdev-eks",
                "auth_method": "eks"
            },
            "terraform_directory": ".happy/terraform/envs/rdev",
            "task_launch_type": "k8s"
        },
        "staging": {
            "aws_profile": "czi-si",
            "k8s": {
                "namespace": "scruffy-staging-happy-happy-env",
                "cluster_id": "hapi-staging-eks",
                "auth_method": "eks"
            },
            "terraform_directory": ".happy/terraform/envs/staging",
            "task_launch_type": "k8s"
        }
    },
    "services": [
        "scruffy"
    ],
    "features": {
        "enable_dynamo_locking": true,
        "enable_happy_api_usage": true,
        "enable_ecr_auto_creation": true
    },
    "stack_defaults": {
        "create_dashboard": false,
        "routing_method": "CONTEXT",
        "services": {
            "scruffy": {
                "build": {
                    "context": ".",
                    "dockerfile": "Dockerfile"
                },
                "health_check_path": "/health/",
                "name": "scruffy",
                "path": "/health/*",
                "platform_architecture": "arm64",
                "port": 3000,
                "priority": 0,
                "service_type": "INTERNAL",
                "success_codes": "200-499"
            }
        },
    }
}
~~~

#### Breakdown:

* `stack_defaults` - used to configure your base set of options for each of your environment stacks. All the settings in `stack_defaults` will
apply to dev, staging, and prod stack configurations with the option to override these settings in each specific environment. This setting
keeps your deployment configuration DRY. Ideally, only prod should be somewhat different in the memory and CPU.
* `features` - feature flags to enable. We highly recommend `enable_dynamo_locking`, `enable_happy_api_usage`, and `enable_ecr_auto_creation`
as these will soon be made normal features of happy.
* `services` - the names of services to deploy. The services names should match the service keys in your docker-compose.yml when performing `happy create/update`. If this field isn't filled out, you might recieve a validation error.
* `environments` - an object of environment configurations. This is mapping of environment to AWS account. Each environment needs the terraform directory to use as its terraform template, the AWS profile to use when deploying this environment and the EKS cluster information of your [happy environment](./deploy_first_env.md)
* `app` - the name of the application. These are namespaced names that separate other stacks that might live in the same cluster

### Terraform Templates

The terraform template folders in .happy/terraform (or wherever you have configured your terraform to live) act as 
[configuration](../stacks/terraform.md) too.
Every invocation of `happy create` or `happy update` will bundle up these files, inject some variables from your configuration and execute
a terraform apply against this configuration. Technically, any terraform can go in these templates, but we recommend 
utilizing the [`happy-stack-eks`](https://github.com/chanzuckerberg/happy/tree/main/terraform/modules/happy-stack-eks) module to do most
of the heavy lifting. It is well tested and contains a lot of features such a `happy config`, E2E encryption, OIDC authentication for stacks, 
all the networking, and more to access your service from wherever you need it.

The `happy-stack-eks` module has a lot of [variables](../stacks/terraform.md#terraform) 
to configure your service. However, rather than reading all the variables, we would suggest playing with the examples as they might be
more educational than trying to read into every configuration option. All of the configuration options have sane defaults and they
might only be needed in the most advanced of cases. The defaults will work for most standard applications and can be interated on after that.