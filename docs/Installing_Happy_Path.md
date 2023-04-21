# Installing Happy Path

## MacOS

```
brew tap chanzuckerberg/tap
brew install happy fogg tfenv
brew install --cask session-manager-plugin
tfenv install # Add a specific version if required
```

## Linux

* Download the appropriate Happy tarball for your architecture from [the Releases Page on Github](https://github.com/chanzuckerberg/happy/releases).
* Decompress the tarball and place the `happy` binary somewhere in your `$PATH`.
* [Install session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
* Install `fogg` per the [instructions](https://github.com/chanzuckerberg/fogg#linux) on Github.
* [Install Terraform Version Manager](https://github.com/tfutils/tfenv)
  * Run `tfenv install` to install the latest version of Terraform. Note that tfenv works by changing symlinks, so you may need to add yourself to a `tfenv` group or similar. Check the group ownership of `/var/lib/tfenv/version`.


## From Source

First, make sure you have the following installed and in your path:

* [Go](https://go.dev/doc/install)
* [Go Releaser](https://goreleaser.com/install/)

### Install Happy
```
git clone https://github.com/chanzuckerberg/happy.git
cd happy/cli
git checkout <latest version tag>

goreleaser build

# You'll want to make sure $HOME/go/bin/ exists and is in your $PATH
cp dist/<arch>/happy $HOME/go/bin/

```

This will build Happy for several architectures under the `dist` directory. Copy the binary from the appropriate arch folder to someplace in your $PATH. The example above uses `/usr/local/bin`, but that's not a requirement.

### Install Fogg

```
git clone https://github.com/chanzuckerberg/happy.git
cd fogg
git checkout <latest version tag>

make install
```

This will install `fogg` to your `$HOME/go/bin` directory. If you'd prefer it somewhere else, feel free to move it anywhere in your $PATH. If you decide to leave it, make sure `$HOME/go/bin` is in your $PATH.

### Install `tfenv`

* [Install Terraform Version Manager](https://github.com/tfutils/tfenv)
* Ensure that your user is added to the group which owns `/var/lib/tfenv/version` if using Linux.
* Run `tfenv install` to install the latest version of Terraform. Optionally specify a version of Terraform to use. 


### Install Session Manager Plugin

You can start by [Follow the instructions for your OS](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html).

If your Linux Distribution is not listed here, you may have more luck following the [build steps using Docker](https://github.com/aws/session-manager-plugin) to build the plugin. Once the docker build completes, you can copy the plugin for your architecture from `bin/<arch>_plugin` to somewhere in your $PATH.

