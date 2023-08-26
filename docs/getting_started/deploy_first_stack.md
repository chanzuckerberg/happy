---
parent: Getting Started
layout: default
has_toc: true
---

# Deploy Your First Stack
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Summary

Generally, a happy environment would need to be set up before you can start deploying stacks--but that is boring. It is much more fun to 
to start using happy and deploy your first application! Luckily, we have a sandbox environment to do just that. 

## Ground Rules

Some ground rules:

* This is a very handy way to play with happy and experiment with new stack ideas--**please don't abuse it**
* Stacks are cleaned up automatically at the end of the week
* This is not a production environment
* The sandbox is for CZI employees only

## Deploy

Here's how to get started:

* Install the latest version of [happy CLI](./installation.md)
* Clone [happy](https://github.com/chanzuckerberg/happy/tree/main)
* Navigate to [./happy/tree/main/examples/typical_app](https://github.com/chanzuckerberg/happy/tree/main/examples/typical_app) (or any of the other example projects in /examples)
* Execute `happy create <stackname>` where `<stackname>` is the name of the stack
* You're done!

This stack will deploy two services. To see the endpoints to access them, use `happy list` and it should look like the following:

~~~
$ happy list
[INFO]: Listing stacks from the happy api
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ NAME     │ OWNER    │ APP     │ REPO                                    │ BRANCH                     │ HASH                           │ STATUS  │ URLS                                                          │ LASTUPDATED │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ typical1 │ alokshin │ typical │ https://github.com/chanzuckerberg/happy │ alokshin/global-name-check │ dirty git tree (PLEASE COMMIT  │ applied │ https://typical1.happy-playground-rdev.rdev.si.czi.technology │ 176h34m39s  │
│          │          │         │                                         │                            │ YOUR CHANGES)                  │         │                                                               │             │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
~~~

Change the Dockerfile, application code, happy configuration however you like and redeploy using `happy update <stackname>`. 

## Inspect

Once your stack has successfully been deployed, you can obviously use your web browser to navigate to the application. But there are some
other cool commands you can experiment with a running stack.

### Logging & Monitoring

* `happy logs <stack_name> <service_name>` - print the logs of the running application. Many of my stacks have way too many logs and I often 
only want to see the logs from a recent error. To do that, I often will do `happy logs <stack_name> <service_name> --since 10m --output /tmp/logs.json`. Then I open the logfile in my editor for further inspection. This command does not stream logs, but happy does capture all logs 
in CloudWatch. At the end of the command it will ask if you'd like to view the full log stream. This will send you to the AWS console where you can see all the application logs.
* `happy events <stack_name>` - prints the latest events related to the stack's deployment. This is handy for when a deployment doesn't go so
well. For instance, you might have specified the wrong port number, a health check is failing, or you accidently are deploying on the wrong
architecture. All these types of events will either show up in `happy events` or `happy logs`.
* `happy resources <stack_name>` - prints the AWS resources created for your stack

~~~
$ happy resources typical1
[INFO]: Retrieving stack 'typical1' from environment 'rdev'
│─────────────────────────────────────────────────────────────────────────│────────────────────────────────────────│─────────────────────────────────────────│───────────│─────────────────────────────────────────────────────────────────────────────────────────│
│ MODULE (50)                                                             │ NAME                                   │ TYPE                                    │ MANAGEDBY │ INSTANCES                                                                               │
│─────────────────────────────────────────────────────────────────────────│────────────────────────────────────────│─────────────────────────────────────────│───────────│─────────────────────────────────────────────────────────────────────────────────────────│
│ module.stack                                                            │ stack_dashboard                        │ datadog_dashboard_json                  │ terraform │ rde-hbv-vbw                                                                             │
│ module.stack                                                            │ oidc_config                            │ kubernetes_secret                       │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-oidc-config                                   │
│ module.stack                                                            │ suffix                                 │ random_pet                              │ terraform │ helping-macaque                                                                         │
│ module.stack.module.services["frontend"]                                │ deployment                             │ kubernetes_deployment_v1                │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend                                      │
│ module.stack.module.services["frontend"]                                │ hpa                                    │ kubernetes_horizontal_pod_autoscaler_v1 │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend                                      │
│ module.stack.module.services["frontend"]                                │ service                                │ kubernetes_service_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend                                      │
│ module.stack.module.services["frontend"]                                │ this                                   │ random_pet                              │ terraform │ clever-possum                                                                           │
│ module.stack.module.services["frontend"].module.ecr                     │ lifecycle                              │ aws_ecr_lifecycle_policy                │ terraform │ typical1/rdev/frontend                                                                  │
│ module.stack.module.services["frontend"].module.ecr                     │ repo                                   │ aws_ecr_repository                      │ terraform │ arn:aws:ecr:us-west-2:401986845158:repository/typical1/rdev/frontend                    │
│ module.stack.module.services["frontend"].module.iam_service_account     │ role                                   │ aws_iam_role                            │ terraform │ []                                                                                      │
│ module.stack.module.services["frontend"].module.iam_service_account     │ service_account                        │ kubernetes_service_account              │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend-rdev-typical1                        │
│ module.stack.module.services["frontend"].module.ingress[0]              │ ingress                                │ kubernetes_ingress_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend                                      │
│ module.stack.module.services["frontend"].module.ingress[0]              │ ingress_bypasses                       │ kubernetes_ingress_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-frontend-options-bypass                       │
│ module.stack.module.services["internal-api"]                            │ deployment                             │ kubernetes_deployment_v1                │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api                                  │
│ module.stack.module.services["internal-api"]                            │ hpa                                    │ kubernetes_horizontal_pod_autoscaler_v1 │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api                                  │
│ module.stack.module.services["internal-api"]                            │ service                                │ kubernetes_service_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api                                  │
│ module.stack.module.services["internal-api"]                            │ this                                   │ random_pet                              │ terraform │ smashing-catfish                                                                        │
│ module.stack.module.services["internal-api"].module.ecr                 │ lifecycle                              │ aws_ecr_lifecycle_policy                │ terraform │ typical1/rdev/internal-api                                                              │
│ module.stack.module.services["internal-api"].module.ecr                 │ repo                                   │ aws_ecr_repository                      │ terraform │ arn:aws:ecr:us-west-2:401986845158:repository/typical1/rdev/internal-api                │
│ module.stack.module.services["internal-api"].module.iam_service_account │ role                                   │ aws_iam_role                            │ terraform │ arn:aws:iam::401986845158:role/si-playground-eks-v2/typical1-internal-api-rdev-typical1 │
│ module.stack.module.services["internal-api"].module.iam_service_account │ service_account                        │ kubernetes_service_account              │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api-rdev-typical1                    │
│ module.stack.module.services["internal-api"].module.ingress[0]          │ ingress                                │ kubernetes_ingress_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api                                  │
│ module.stack.module.services["internal-api"].module.ingress[0]          │ ingress_bypasses                       │ kubernetes_ingress_v1                   │ terraform │ si-rdev-happy-eks-rdev-happy-env/typical1-internal-api-options-bypass                   │
│                                                                         │ *typical1-frontend-rdev-typical1       │ ServiceAccount                          │ k8s       │                                                                                         │
│                                                                         │ *typical1-internal-api-rdev-typical1   │ ServiceAccount                          │ k8s       │                                                                                         │
│                                                                         │ typical1-frontend                      │ Service                                 │ k8s       │                                                                                         │
│                                                                         │ typical1-internal-api                  │ Service                                 │ k8s       │                                                                                         │
│                                                                         │ *typical1-oidc-config                  │ Secret                                  │ k8s       │                                                                                         │
│                                                                         │ typical1-frontend                      │ Deployment                              │ k8s       │                                                                                         │
│                                                                         │ typical1-internal-api                  │ Deployment                              │ k8s       │                                                                                         │
│                                                                         │ typical1-frontend                      │ HorizontalPodAutoscaler                 │ k8s       │                                                                                         │
│                                                                         │ typical1-internal-api                  │ HorizontalPodAutoscaler                 │ k8s       │                                                                                         │
│                                                                         │ typical1-frontend                      │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-frontend                                           │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-frontend                                           │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│                                                                         │ typical1-frontend-options-bypass       │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-frontend-options-bypass                            │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-frontend-options-bypass                            │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│                                                                         │ typical1-internal-api                  │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-internal-api                                       │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-internal-api                                       │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│                                                                         │ typical1-internal-api-mybypass-bypass  │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-internal-api-mybypass-bypass                       │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-internal-api-mybypass-bypass                       │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│                                                                         │ typical1-internal-api-mybypass2-bypass │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-internal-api-mybypass2-bypass                      │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-internal-api-mybypass2-bypass                      │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│                                                                         │ typical1-internal-api-options-bypass   │ Ingress                                 │ k8s       │                                                                                         │
│ k8s:Ingress/typical1-internal-api-options-bypass                        │                                        │ Application Load Balancer               │ k8s       │ k8s-stacktypical1help-2fe1caf63c-78548467.us-west-2.elb.amazonaws.com                   │
│ k8s:Ingress/typical1-internal-api-options-bypass                        │                                        │ Route 53 Entry                          │ k8s       │ typical1.happy-playground-rdev.rdev.si.czi.technology                                   │
│─────────────────────────────────────────────────────────────────────────│────────────────────────────────────────│─────────────────────────────────────────│───────────│─────────────────────────────────────────────────────────────────────────────────────────│
~~~
* `happy get <stack_name>` - print metadata information about the deployed stack

~~~
$ happy get typical1
[INFO]: Retrieving stack 'typical1' from environment 'rdev'
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ NAME     │ OWNER    │ APP     │ REPO                                    │ BRANCH                     │ HASH                           │ STATUS  │ URLS                                                          │ LASTUPDATED │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ typical1 │ alokshin │ typical │ https://github.com/chanzuckerberg/happy │ alokshin/global-name-check │ dirty git tree (PLEASE COMMIT  │ applied │ https://typical1.happy-playground-rdev.rdev.si.czi.technology │ 177h2m8s    │
│          │          │         │                                         │                            │ YOUR CHANGES)                  │         │                                                               │             │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│─────────────────────────│──────────────────────────────────────────────────────────────────────────────│
│ RESOURCE (23)           │ VALUE                                                                        │
│─────────────────────────│──────────────────────────────────────────────────────────────────────────────│
│ Environment             │                                                                              │
│ TFE                     │                                                                              │
│   Environment Workspace │ https://si.prod.tfe.czi.technology/app/happy-playground/workspaces/env-      │
│   Stack Workspace       │ https://si.prod.tfe.czi.technology/app/happy-playground/workspaces/-typical1 │
│   Backlog size          │ 2 outstanding runs                                                           │
│                         │ edu-platform-infra                                                           │
│                         │ (planning)->1                                                                │
│                         │ lp-infra (applying)->1                                                       │
│ AWS                     │                                                                              │
│   Account ID            │ [401986845158]                                                               │
│   Region                │ us-west-2                                                                    │
│   Profile               │ czi-playground                                                               │
│ Service                 │ frontend                                                                     │
│   Compute               │ K8S                                                                          │
│   namespace             │ si-rdev-happy-eks-rdev-happy-env                                             │
│   deployment_name       │ typical1-frontend                                                            │
│   auth_method           │ eks                                                                          │
│   kube_api              │ https://FF38057B0DB0747C32795D4B8265796A.gr7.us-west-2.eks.amazonaws.com     │
│ Service                 │ internal-api                                                                 │
│   Compute               │ K8S                                                                          │
│   namespace             │ si-rdev-happy-eks-rdev-happy-env                                             │
│   deployment_name       │ typical1-internal-api                                                        │
│   auth_method           │ eks                                                                          │
│   kube_api              │ https://FF38057B0DB0747C32795D4B8265796A.gr7.us-west-2.eks.amazonaws.com     │
│─────────────────────────│──────────────────────────────────────────────────────────────────────────────│
~~~


### Iterate and Develop

* `happy shell <stack_name> <service_name>` - shell into your application. This is very handy for debugging running applications. The shell
interactions are authenticated using IAM permissions.
* `happy config <cp,delete,diff,exec,get,list,set>` - add a key, value pair to your deployment. A key, value pair will be stored and added 
to your running container as an environment variable. Key, value pairs can be set on the stack-level or on the environment-level 
(all stacks in the environment get the same key value pairs).
* `happy restart <stack_name>` - performs a quick redeployment of the running containers. This is handy if you want to update your `happy config` of your 
application and have your application pick up the changes. The docker images will be the same and no new code will be uploaded. Everything will be the same
except your containerized applications will be redeployed with updated configuration values.
* `happy delete <stack_name>` - deletes the stack
