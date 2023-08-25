---
parent: Development
title: Deployment
layout: default
---

# Deployment

## Structure 

SASTisfaction is deployed as a docker container inside an ECS cluster. We have two environments living inside of the czi-sec AWS account. The flow of deployment is as follows:

* For every push, Github CI will create a container of that code. All containers are published to AWS ECR in czi-sec
* For every push to the main branch, Github CI will deploy the latest ECR container to staging and force a new cluster deployment
* For every published release, Github CI will deploy the latest ECR container to production and force a new cluster deployment. 
* New cluster deployments spin up new instances of our containers. Once the containers are spun up, the old container are taken down.
* All containers are deployed with Fargate, so we do not manage the underlying AWS instances. Follow this deployment flow and we should have little to worry about with maintaining the network infrastructure.

## Releases

We use [release-please](https://github.com/googleapis/release-please) to do all our releases. This means to create a release, [go to the PRs section](https://github.com/chanzuckerberg/sastisfaction/pulls), approve, and merge the release PR. The Github Action will handle the rest.

To make sure that release-please collects all the previous PR information properly, make sure all your PRs are in the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) format.