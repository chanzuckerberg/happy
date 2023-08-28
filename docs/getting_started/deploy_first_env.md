---
parent: Getting Started
layout: default
has_toc: true
---

# Deploy You First Happy Environment

## Background

Happy doesn't really care how you build the [long-lived infrastructure](../design.md#environment). At CZI, we have a very specific 
way of doing this outlined below, but it isn't required. The reason it is not required is because happy defines a contract 
that should be agreed upon between the stacks and the happy environment. This contract is called the
[integration-secret](../design.md#integration-secret).
As long as an integration secret is implemented and fulfills the expected values, happy can deploy to it. This allows users with 
existing infrastructure to easily port their environments to happy. Simply create a JSON document and upload it
as a K8S opaque secret called integration-secret and your stacks will be able to find all the infra it needs to deploy.

For more green field projects, we recommend the CZI approach below.

## CZI Happy Environment Terraform Module

CZI has serveral prebuild modules for building a happy environment. They are listed here in the order they should be executed:

* [route-53](https://github.com/chanzuckerberg/happy/tree/main/terraform/modules/happy-route53)
  * domain for your stacks
* [aws-env](https://github.com/chanzuckerberg/shared-infra/tree/main/terraform/modules/aws-env)
  * VPC
* [eks-cluster-v2](https://github.com/chanzuckerberg/shared-infra/tree/main/terraform/modules/eks-cluster-v2)
  * EKS cluster
* [k8s-core](https://github.com/chanzuckerberg/shared-infra/tree/main/terraform/modules/k8s-core)
  * EKS cluster add-ons
* [happy-env-eks](https://github.com/chanzuckerberg/happy/tree/main/terraform/modules/happy-env-eks)
  * K8S namespace and integration-secret

Each happy environment needs a AWS route53 zone ID, VPC, EKS cluster, and integration secret. The integration secret should
have all the fields populated by the above modules. See the [integration-secret docs](../design.md#integration-secret) on 
what fieldsa are required to be filled in. If you are using the above modules, you don't need to worry about this. It will create
all the necessary long-lived infrastructure and create an integration secret in the proper place.

Ideally, we'd like to each of these modules in a separate Terraform Workspace and be mapped to a specific environment. Here's an
example directory structure we usually see for happy projects:

~~~
└── terraform
    ├── envs
    │   ├── prod
    │   │   ├── cloud-env
    │   │   ├── eks
    │   │   ├── happy-eks
    │   │   ├── k8s-core
    │   │   ├── route53
    │   │   └── secrets_from_aws_param
    │   ├── rdev
    │   │   ├── cloud-env
    │   │   ├── eks
    │   │   ├── happy-eks
    │   │   ├── k8s-core
    │   │   ├── route53
    │   │   └── secrets_from_aws_param
    │   └── staging
    │       ├── cloud-env
    │       ├── eks
    │       ├── happy-eks
    │       ├── k8s-core
    │       ├── route53
    │       └── secrets_from_aws_param
~~~

If you need more advanced infrastructure components that don't come with these modules, add components under each
of the environments for those pieces. Try not to add these new terraform components to the existing modules as it will 
make them bloated and more likely to fail during a terraform apply.

## Adding Custom Infra to Happy Environment

Happy comes with some things out of the box:

* Aurora postgres
* S3 buckets
* Batch

However, there might be other types of AWS resources to deploy the application (ie CloudFront, Redis, SQS, SNS, etc.). Add these
terraform elements as you normally would to your environment. Once applied, we need to add the corresponding elements to the integration 
secret value. This will allow your running stacks to see the infrastructure and connect to it. To add your custom infrastrucutre to your
environment, add the values you want your stacks to consume at the [additional_secrets](https://github.com/chanzuckerberg/happy/blob/fcb0fad658ee0cecd01921dd0cb3f45901cfaf68/terraform/modules/happy-env-eks/variables.tf#L52) section of the 
[happy-env-eks](https://github.com/chanzuckerberg/happy/tree/main/terraform/modules/happy-env-eks) module.

