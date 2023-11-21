# Happy Path Deployment Tool

Visit the Happy Path documentation for more details: https://chanzuckerberg.github.io/happy/

The Happy Path Deployment Tool is an open-source project led by the Chan Zuckerberg Initiative (CZI). It is a platform for deploying and managing containerized applications at scale in adherence to CZI security practices. The tool is designed to be easy to use and operate, and it is available for both on-premises and cloud deployments. Happy builds and deploys your application, and once it is released, helps you support it.

Happy Path is based on these principles:

* Operational simplicity: Happy Path takes the bite out of a complex container orchestration operations and infrastructure management
* Extensibility: Happy Path functionality can be extended through various Terraform hooks, custom Terraform modules and Helm Chart customization
* Reliability: Happy Path is reliable and is production-ready, used by multiple engineering teams at CZI

### Security

Please note: If you believe you have found a security issue, please responsibly disclose by contacting us at security@chanzuckerberg.com


### Repository structure

This project is a monorepo for the components of a Happy ecosystem:
* `cli/` - the `happy` CLI tool
* `api/` - the `happy` API server
* `shared/` - components shared between `api`, `cli` and the `terraform/provider`
* `hvm/` - the `happy` version manager
* `terraform/` - a collection of TF modules we use to provision long-lived infrastructure and application stacks
* `terraform/provider` - `happy` terraform provider
* `helm-charts/stack` - an experimental helm chart
* `examples/` - sample applications that illustrate various concepts that `happy` supports, such as sidecars, tasks, multi-service deployments, GPU intense workloads, and so on


### Features
* Manages short-lived infrastructure (we deploy into your compute)
* Groups services together (we call it a `stack` for co-deployment), each `stack` is isolated, and you can have multiple `stacks` created for the same application.
* Easily promote changes from lower to higher environments
* Supports automated deployments through github workflows
* Has an extensive set of Github workflows
* Supports both AWS ECS and EKS runtimes, and allows for an easy migration between the two
* Abstracts out IaC code with the intent that developers should only focus on the application code
* Supports Linkerd service mesh, mTLS and service-to-service authorization when deployed on EKS with Linkerd installed
* Plays nicely with `external-dns`, `karpenter`, `cluster-autoscaler`
* Integrates with Datadog for dashboarding (assuming you have a `datadog` agent deployed into your EKS)
* Provides service discovery and autoscaling capabilities
* Supports both amd64 and arm64 architectures
* Supports metrics collection, health monitoring through healthchecks, and synthetics
* Supports rolling updates to minimize downtime

### Prerequisites

You will need to have Docker desktop, AWS CLI, and `terraform` installed to use Happy. You can install them and other useful tools by running 

```
brew tap chanzuckerberg/tap
brew install awscli helm kubectx kubernetes-cli aws-oidc linkerd jq terraform
brew install --cask docker
```

Docker Desktop needs to be running; and aws cli needs to be configured by running `aws configure`, with profiles setup. Make sure `cat ~/.aws/config` does not return an empty string (we assume you already have an AWS account).

In addition to the above, you will need an up and running EKS cluster, that contains a happy environment namespace (it contains a secret called `integration-secret`).

### Install

Install `happy`:

#### MacOS
```sh
brew tap chanzuckerberg/tap
brew install happy
```

Alternatively, you can install hvm and install `happy` using [hvm](https://github.com/chanzuckerberg/happy/blob/main/hvm/README.md):
```
brew tap chanzuckerberg/tap
brew install hvm
hvm install chanzuckerberg happy <version>
hvm set-default chanzuckerberg happy <version>
```

#### Linux

Binaries are available on the releases page. Download one for your architecture, put it in your path and make it executable.

Instructions on downloading the binary:

1. Go here: <https://github.com/chanzuckerberg/happy/releases> to find which version of happy you want.
2. Run `curl -s https://raw.githubusercontent.com/chanzuckerberg/happy/master/download.sh | bash -s -- -b HAPPY_PATH VERSION`
   1. HAPPY_PATH is the directory where you want to install happy
   2. VERSION is the release you want
3. To verify you installed the desired version, you can run `happy version`.


### Getting started

#### Setting up a brand new application
Create a folder called `myapp`
```sh
mkdir myapp
cd myap
```

Bootstrap the Happy application:

```sh
happy bootstrap --force
```

Answer the prompts:

* `What would you like to name this application?`: `myapp`. 
* `Your application will be deployed to multiple environments. Which environments would you like to deploy to?`: `rdev`
* `Which aws profile do you want to use in rdev?`: select the appropriate aws configuration profile
* `Which aws region should we use in rdev?`: select the aws region with the EKS cluster
* `Which EKS cluster should we use in rdev?`: select the cluster you will be deploying to
* `Which happy namespace should we use in rdev?`: select the namespace that has a Happy environment configured
* `Would you like to use dockerfile ./Dockerfile as a service in your stack?`: `Y`
* `What would you like to name the service for ./Dockerfile?`: `myapp`
* `What kind of service is myapp?`: `EXTERNAL`
* `Which port does service myapp listen on?`: use a port number other than 80 or 443
*  `Which uri does myapp respond on?`: `/`
* `File /tmp/myapp/docker-compose.yml already exists. Would you like to overwrite it, save a backup, or skip the change?`: overwrite (if prompted)

At this point, your folder structure looks like
```
.
├── .happy
│   ├── config.json
│   └── terraform
│       └── envs
│           └── rdev
│               ├── main.tf
│               ├── outputs.tf
│               ├── providers.tf
│               ├── variables.tf
│               └── versions.tf
├── Dockerfile
└── docker-compose.yml
```

Happy configuration is blended from three sources: `config.json` for environment and application structure setup; `main.tf` to wire the terraform code and provide baseline parameters, and `docker-compose.yaml` to indicate where relevant `Dockerfile` files are located. Multiple environments (think `dev`, `staging`, and `prod`) can be defined with unique configuration settings.

Let's create a new stack: `happy create myapp-stack`, not that we have namespaced the stack name, as stack names are unique in the entire environment due to DNS constraints. Once the stack is created, `happy` will display a list of endpoints:
```
[INFO]: service_endpoints: {
	"EXTERNAL_MYAPP_ENDPOINT": "https://myapp-stack.<PARENT-DNS-ZONE-NAME>",
	"PRIVATE_MYAPP_ENDPOINT": "http://myapp-myapp.<NAMESPACE>.svc.cluster.local:80"
}
```

`EXTERNAL_MYAPP_ENDPOINT` is accessible from your browser. `PRIVATE_MYAPP_ENDPOINT` is accessible by other applications on running on the same cluster (if Linkerd is enabled, you can have granular controls over who can connect to it). Try curl-ing the `EXTERNAL_MYAPP_ENDPOINT`. You will get a default nginx response.

Now, list stacks out: `happy list`, you will get a human readable output. For machine-readable output, run `happy list --output json`. If you intend to update the application after changes were made, run `happy update myapp-stack`. Close the session by deleting the stack: `happy update myapp-stack`.


#### Sample apps
Clone this repo: 
```sh
git clone https://github.com/chanzuckerberg/happy.git
```

Navigate to an example app and try happy out:
```sh
cd examples/typical_app
happy list
happy create mystack
happy update mystack
happy delete mystack
```

### Contributing

This project adheres to the Contributor Covenant code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to opensource@chanzuckerberg.com.
//
