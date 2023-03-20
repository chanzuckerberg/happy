---
title: Getting Started with Happy
---
Getting Started with Happy
This means that all stack-level resources need to be named so that they don’t conflict with any other stack, but also so that they don’t conflict with other stacks in other Happy Path environments!

AWS namespacing rules make this even more complicated - for example, ECS services are namespaced to the ECS cluster they belong to, but ECS task definitions are namespaced at the AWS account level!

### General naming
Where resources are created in the AWS account-wide namespace, the recommended resource naming scheme is:


`${app_name}-${env_name}-${stack_name}-${resource_name}`

Where resources are created in deeper namespaces, as with services in an ECS cluster, it’s appropriate to drop any prefixes that aren’t required to prevent collisions.

### ECR naming
Since ECR repositories are intended to be used across the entire application development lifecycle (dev, staging, and production), the recommended naming scheme for ECR is:


`${app_name}-${image_label}`

For example, the repository that manages the backend service image for the widgetco application would be widgetco-backend
