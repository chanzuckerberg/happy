---
parent: Getting Started
layout: default
has_toc: true
---

# Installing Happy Path
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Prereqs

Before getting started, happy relies on the following: `happy` CLI tool, [`terraform`](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli), and `aws` [CLI tool](https://formulae.brew.sh/formula/awscli).

### MacOS

```
brew tap chanzuckerberg/tap
brew install happy fogg tfenv awscli
tfenv install
```

### Linux

* Download the appropriate Happy tarball for your architecture from [the Releases Page on Github](https://github.com/chanzuckerberg/happy/releases).
* Decompress the tarball and place the `happy` binary somewhere in your `$PATH`.
* [Install Terraform](https://github.com/tfutils/tfenv)
  * Run `tfenv install` to install the latest version of Terraform. Note that tfenv works by changing symlinks, so you may need to add yourself to a `tfenv` group or similar. Check the group ownership of `/var/lib/tfenv/version`.

