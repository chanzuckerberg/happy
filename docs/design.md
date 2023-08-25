---
layout: default
title: Design
nav_order: 3
---
 
# Design
 
There are two main pieces of happy's design:

1. Long-lived infrastructure (aka "happy environments")
1. Short-lived infrastructure (aka "happy stacks")

<img width="long-lived-short-lived" alt="Screenshot 2023-08-25 at 4 07 51 PM" src="https://github.com/chanzuckerberg/happy/assets/76011913/333cbfca-8b0e-40f8-84a5-5a49be3d69a1">

## Environment

The long-lived infrastructure is what we generally call the happy environment. It is what you deploy to. It includes the following things:

* EKS or ECS cluster
* Networking
* CI/CD workflows
* Batch processing
* Aurora Postgres DB cluster
* OIDC 
* S3 buckets
* Terraform workspaces
* Any other custom or advanced infra

We refer to it as long-lived instructure because once it is created, it does not change a lot. These are generally resources we don't want
to bring up and down a lot or take a long time to provision. They might have data that can't be destroyed. The long-lived infrastructure
is what we build our stacks on top of.

## Stacks

The short-lived infrastructure is what we generally call happy stacks. Each stack represents a fully encapsulated version of your application
to deploy. Stacks are:

* A containerized application
* Containerized jobs for batch, cron or other processing
* Runtime configuration (key, value pairs)
* Any small custom infrastructure

Each stack is a complete running application. They are meant to come up and down very quickly. They are great for:

* Testing a new feature on a PR
* Showcasing your work in a demo environment
* Spinning up a pentesting environment
* Spinning up a development environment
* Iteration on microservice locally

## Integration Secret

Happy stacks need to know where to deploy to. This information is conveyed in a JSON documented called the integration secret. The documents
looks like this:

~~~json
{
    "certificate_arn": "arn:aws:acm:us-west-2:401986845158:certificate/fe87ad86-23b0-41b4-b2ff-e1724a004d59",
    "ci_roles": [
        {
            "arn": "arn:aws:iam::401986845158:role/gh_actions_si_playground_eks_v2",
            "name": "gh_actions_si_playground_eks_v2"
        }
    ],
    "cloud_env": {
        "database_subnet_group": "si-playground",
        "database_subnets": [
            "subnet-00fc2b9327edbc4d9",
            "subnet-0fb4dd8fa2a6e5f18",
            "subnet-0a8964bf1c4d576e4",
            "subnet-0b8efc311727e0255"
        ],
        "private_subnets": [
            "subnet-0ffdf3a2f10adb1f9",
            "subnet-0d69f41508803494c",
            "subnet-08c029e8a4a6b6a32",
            "subnet-0fa1ce01ef263099c"
        ],
        "public_subnets": [
            "subnet-0bfdc8898996f59f2",
            "subnet-0af8d4c628578cd67",
            "subnet-0aa5f4235faa1c475",
            "subnet-0881e4e79011778cc"
        ],
        "vpc_cidr_block": "10.252.0.0/16",
        "vpc_id": "vpc-00c10bee0068a11f2"
    },
    "dbs": {},
    "dynamo_locktable_name": "si-prod-happy-eks-prod-stacklist",
    "ecrs": {},
    "eks_cluster": {
        "cluster_arn": "arn:aws:eks:us-west-2:401986845158:cluster/si-playground-eks-v2",
        "cluster_ca": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1EVXlOREl3TWpnME1sb1hEVE16TURVeU1USXdNamcwTWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTjNYCjduRzRhREQxcHRpb2Q1enJjNlFtZ1YyMnZtVWk5azFWY0tjRDE3WDAzOUk3dUc5S3JXRURxUXlrSGpaYmg0VzQKYzN6YmtFZW9FMlJ0R1p2TGowL1duLzNuazl3ZElETnZUb1pQd3ZqbGtFSVVQVnFrMEc3a3hYREt6M3VnSVBRMgpXQnlHOGJLU0FpUDlneXp3dktYTWVjSkEvYVBpVUhZMkhWRFFPbG5SNkM0dTFqMG9md2I5di9ReHY4Y3NPbE1ECjJIUFh4TVBSODdTZ2UxQzRvelZMTGNzS0Vyc3hYZ0NNd1h6T21xajBHTjg4RVBJSnVFQUQwcHF4aU9nRi90YXUKek92S1VVcTFUMGFCektqVjVqVWZGYXU4aGUyTHBYZzI4cUFKbzlpZGNibXk3WG5PTGhrdjgzdDlRcEVlYjRVegpRcldrZ1ZINzN1Q0g3Q2NoTExVQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZLQkdvZWg2di9EejM5WjNVbDFPemJQMkl4Z1BNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSzRCOEZrSVp5cW15N1A3TlZKQQpxNzRoL3lieXpKcUZWU3VrdmVCd0tkSlE1WnEyaU1UeldaQU95K2pqY2xtZkYrMDJBU0xyMWlSRnYwaXV5MEROCjh5V0pja3ZNcWJRTm9raFhoVFltMjBHbmVNaEd6cW9jNUlmTmpKTUlVQzlaSFUzYVRtV3ZQVXFVM0RKd0FQZHQKT0pWMTR3Y3F4WWxMMkJ6dElwOUlpK0ZyNmp5NzFWT25haEphV1VudHNjMnFpTDZ1cE5ZMjk1TmkwV0R2V00vVQoyWmdYTHRNemVvanMxdWdIUWZkb3NuRGo2OUZFUDk4V1VrQzg5NjVxcno4SjV3ZENvVnI2ZnMvWGRPeElDOWtwCnNzVHlldm1Jb1pKalBiQmJNQnhEWHVaMHBSUmxpbnZRd3A0RElPbUdRN0w2WmQvd1pGa2FWOFRxbXdpQXJhdDIKM0lrPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
        "cluster_endpoint": "https://FF38057B0DB0747C32795D4B8265796A.gr7.us-west-2.eks.amazonaws.com",
        "cluster_id": "si-playground-eks-v2",
        "cluster_oidc_issuer_url": "https://oidc.eks.us-west-2.amazonaws.com/id/FF38057B0DB0747C32795D4B8265796A",
        "cluster_version": "1.26",
        "oidc_provider_arn": "arn:aws:iam::401986845158:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/FF38057B0DB0747C32795D4B8265796A",
        "worker_iam_role_name": "si-playground-eks-v2-eks-node20230524200454767700000001",
        "worker_security_group": "sg-0d05e6512e80e6a75"
    },
    "external_zone_name": "happy-playground-prod.prod.si.czi.technology",
    "hapi_config": {
        "assume_role_arn": "arn:aws:iam::401986845158:role/tfe-si",
        "base_url": "https://hapi.hapi.prod.si.czi.technology",
        "kms_key_id": "358d3b31-aa91-4017-b266-214001e37d41",
        "oidc_authz_id": "aus8zrryrcU8fOk9y5d7",
        "oidc_issuer": "0oa8zrudz1DLeTqOG5d7",
        "scope": "happy-eks-prod"
    },
    "kind": "k8s",
    "oidc_config": {
        "client_id": "0oa8zrsvuh6c9hoHe5d7",
        "client_secret": "KAzUaIlYMqHHkbK2AnF1xXhuCsD1j84GL06WwxGg",
        "config_uri": "https://0oa8zrsvuh6c9hoHe5d7:KAzUaIlYMqHHkbK2AnF1xXhuCsD1j84GL06WwxGg@czi-prod.okta.com/oauth2/",
        "idp_url": "czi-prod.okta.com"
    },
    "tags": {
        "env": "prod",
        "managedBy": "terraform",
        "owner": "infra-eng@chanzuckerberg.com",
        "project": "si",
        "service": "happy-eks-prod"
    },
    "tfe": {
        "org": "happy-playground",
        "url": "https://si.prod.tfe.czi.technology"
    },
    "vpc_id": "vpc-00c10bee0068a11f2",
    "zone_id": "Z10432202W0G3BIMH0531"
}
~~~

This is a lot of information in here but the general outline looks like the following:

* VPC information
* Database information (username, host, password)
* EKS/ECS cluster to deploy to 
* Tags
* OIDC information (for protected environments)
* TLS certificate information
* IAM roles 

When a stack is deploying, the stack will have access to the integration secret so that it can deploy in the correct place. It 
can also be used at runtime too. When a stack is deployed, the integration secret will be mounted to the deployed container
at `/var/happy/integration-secret`. Applications can use this to connect to databases, assume roles or parse other advanced
infrastructure that might be needed at runtime by the application.

# References:

* [Presentation](https://docs.google.com/presentation/d/1zgbTF_1oq96npmKXxHKFVn5rO96wEsQlj5bgd7axLNA/edit#slide=id.p)
* [Stack options](https://github.com/chanzuckerberg/happy/blob/0142e747802df4768f1d27e27a062f86a821316d/terraform/modules/happy-stack-eks/variables.tf#L46)