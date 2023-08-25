---
layout: default
title: Design
nav_order: 3
---
 
# Design
 
There are two main pieces of happy's design:

1. Long-lived infrastructure (aka "happy environments")
1. Short-lived infrastructure (aka "happy stacks")

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

# References:

* [Presentation](https://docs.google.com/presentation/d/1zgbTF_1oq96npmKXxHKFVn5rO96wEsQlj5bgd7axLNA/edit#slide=id.p)
* [Stack options](https://github.com/chanzuckerberg/happy/blob/0142e747802df4768f1d27e27a062f86a821316d/terraform/modules/happy-stack-eks/variables.tf#L46)