# Installing Happy Path

## MacOS

```
brew tap chanzuckerberg/tap
brew install happy fogg tfenv
tfenv install
```

## Linux

* Download the appropriate Happy tarball for your architecture from [the Releases Page on Github](https://github.com/chanzuckerberg/happy/releases).
* Decompress the tarball and place the `happy` binary somewhere in your `$PATH`.
* Install `fogg` per the [instructions](https://github.com/chanzuckerberg/fogg#linux) on Github.
* [Install Terraform Version Manager](https://github.com/tfutils/tfenv)
  * Run `tfenv install` to install the latest version of Terraform


## From Source

First, make sure you have the following installed and in your path:

* [Go](https://go.dev/doc/install)
* [Go Releaser](https://goreleaser.com/install/)

### Install Happy
```
git clone https://github.com/chanzuckerberg/happy.git
cd happy/cli
git checkout <latest version tag>

export GITHUB_TOKEN=<a Github Personal Access Token>

goreleaser build

cp dist/<arch>/happy /usr/local/bin

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
* Run `tfenv install` to install the latest version of Terraform

