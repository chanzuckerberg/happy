## Pre-requisites

Happy CLI relies on a number of dependencies, listed here:

* Homebrew: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
* AWS CLI: `brew install awscli`
* SessionManager plugin: `brew install --cask session-manager-plugin`
* Docker engine, docker daemon needs to be running
* Docker-compose v2 (in Docker desktop preferences, make sure “Use Docker Compose v2” is enabled - in the image below it is not checked)
* Terraform installed and configured (make sure you have access to TFE)
  * brew install terraform
  * terraform login si.prod.tfe.czi.technology
* AWS-OIDC: 

~~~
mkdir ~/.aws
brew tap chanzuckerberg/tap
brew install aws-oidc
aws-oidc configure --issuer-url https://czi-prod.okta.com --client-id aws-config --config-url https://aws-config-generation.prod.si.czi.technology
~~~

## Installation

### Homebrew

brew install chanzuckerberg/tap/happy

You might run into Error: `chanzuckerberg/tap/happy: Calling bottle :unneeded is disabled! There is no replacement.` To solve this, run brew update-reset

### Direct Download

Download a binary specific for your platform from https://github.com/chanzuckerberg/happy/releases, and place it somewhere in your PATH. For example:

~~~
wget https://github.com/chanzuckerberg/happy/releases/download/v0.4.1/happy_0.4.1_darwin_amd64.tar.gz
tar -xf happy_0.4.1_darwin_amd64.tar.gz
chmod +x ./happy
~~~

Verify installation by running happy version. Then move happy to somewhere in your $PATH. This is sufficient enough to get started with the 
Happy CLI. Run happy list to verify your setup.

## Basic usage

### happy list

Displays a list of all existing stacks in your setup.

### happy build

Builds all docker images referenced in your docker-compose.yml.

### happy push

Builds and pushes docker images to a docker registry (ECR), and tags them appropriately. Default tags look like elopez-chanzuckerberg.com-2022-02-16T21-17-01.

### happy create

Creates a stack using a provided name, for example `happy create stack1`. It will trigger docker build, docker push, and a creation of a new workspace in TFE.

If your happy create errored out, to troubleshoot the issue look into TFE run logs.

### happy update

Builds, pushes docker images and updates an existing stack.

### happy delete

Deletes the entire stack, by running a terraform destroy against the workspace. Don’t delete stacks that errored out on happy create, fix the create and then delete.

### happy migrate

This command migrates the stack (by running the migrate command). 

### happy logs

To get the logs from a running container, execute happy logs STACK_NAME SERVICE_NAME.

### happy hosts

When doing local development (running the stack with docker-compose up), it is useful to create local host name aliases for your own machine. To add/remove host entries associated with your services into /etc/hosts, run sudo happy hosts install or sudo happy hosts uninstall. When host entries are created in /etc/hosts, they look like

# ==== Happy docker-compose DNS for happy-deploy ===
127.0.0.1	frontend.happynet.localdev
# ==== Happy docker-compose DNS for happy-deploy ===

### happy shell

This command allows you to exec into a running container, and requires a Session Manager plugin to be installed.

## Known issues / FAQ

* After several hours, TFE tokens become inactive, and need to be re-activated by logging into TFE. Navigate to https://go.czi.team/tfe . Usually manifests as `failed to get workspace: could not read workspace *******: unauthorized`
* When a stack fails to create, deletion might result in Route53 entries or Cloudwatch log groups left behind, that will fail further creation of the same stack. Don’t happy delete, delete the offending resource and re-plan.
* Docker-compose profiles have limited support
* Docker login will fail if the docker daemon is not running.
* happy create is slow: Run happy push before running happy create or happy update to cache all the docker image layers before hand.
