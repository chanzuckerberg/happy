# Happy Path Deployment Tool


The Happy Path Deployment Tool is an open-source project led by the Chan Zuckerberg Initiative (CZI). It is a platform for deploying and managing containerized applications at scale in adherence to CZI security practices. The tool is designed to be easy to use and operate, and it is available for both on-premises and cloud deployments. Happy builds and deploys your application, and once it is released, helps you support it.

Happy Path is based on these principles:

* Operational simplicity: Happy Path takes the bite out of a complex container orchestration operations and infrastructure management
* Extensibility: Happy Path functionality can be extended through various Terraform hooks, custom Terraform modules and Helm Chart customization
* Reliability: Happy Path is reliable and is production-ready, used by multiple engineering teams at CZI

### Security

Please note: If you believe you have found a security issue, please responsibly disclose by contacting us at security@chanzuckerberg.com

Visit the Happy Path documentation for more details: https://chanzuckerberg.github.io/happy/

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

You will need to have Docker desktop, AWS CLI, and `terraform` installed to use Happy.

### Install

Install `happy`:

#### MacOS
```sh
brew tap chanzuckerberg/tap
brew install happy
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
