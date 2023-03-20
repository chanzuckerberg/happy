---
title: Terminology
---

Terminology:

Stacks:

A "stack" is a set of applications that comprise a particular project or service. For example, the combination of frontend + backend + processing pipelines together can comprise a stack. The shape of a stack is defined as terraform code alongside the applications it's responsible for deploying. Each stack is independent of other stacks and can be created or destroyed without affecting the other stacks.

Happy Environment:

A Happy Path environment is a set of long-lived infrastructure (IAM roles, DB instances, Load balancers, S3 buckets, ECS/EKS clusters, etc) that provide the underlying facilities for applications to be deployed to. The happy-env-ecs terraform module forms the foundation of a Happy Environment. Multiple stacks can be deployed to a Happy Path Environment – for example, all remote-dev stacks live in the same Happy env. The resources in an environment are shared by the stacks, but the stacks maintain their independence from one another.

Integration Secret:

Each Happy Path environment exposes an integration secret in AWS Secrets Manager with information about that environment - ARN's of databases, s3 buckets, networking information, etc. This secret is consumed by the terraform code for application Stacks to configure applications

Happy CLI

The Happy Path team maintains the Happy CLI, which is responsible for managing application Stacks in Happy Environments
