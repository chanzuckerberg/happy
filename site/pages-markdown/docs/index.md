---
title: Happy Path Documentation
---

# Happy Path Documentation and usage

## Getting started

_Ok, but what is it, really?_

Happy Path is a toolkit and a set of recommended best practices for standardizing application lifecycle from local development, remote development, to building, testing, and promotion to production at CZI. Core functionalities:

Empower engineers to develop their applications locally as much as possible for fast iteration

Provide tools to enable engineers to share experimental changes with the rest of their team in a production-like environment

Make it easy to automate deployments to staging and production

Built-in monitoring

 
Where possible, we strongly encourage:

Enable local development without an internet connection (via docker-compose)

Build artifacts [docker images] once, and deploy the same artifact to staging and production

Automatically deploy code to a shared staging environment on every commit to the default branch

Deploy to production as frequently as possible. Yes, automated testing will need to improve!

-

### Yeah yeah, enough of your marketing speak, I want to know about the nuts and bolts!

Really, the core of Happy Path is:
* A set of Terraform Enterprise automation tooling that can dynamically provision new TFE workspaces and apply ad-hoc TF code changes to those workspaces
* A CLI that glues TFE together with docker image management – the CLI handles building images locally and pushing them to a remote repository, and making the resulting tagged images available to Terraform code that references them
* More information about: Happy Path architecture
* A reference implementation of Happy Path infrastructure via terraform modules. Happy Path doesn’t require using these modules, but they are the ‘batteries included’ option.
* Reference implementations of automated application deployment and staging/prod promotion pipelines via GitHub actions