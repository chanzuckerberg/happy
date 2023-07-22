# Happy enVironment Manager

The `H`appy en`V`ironment `M`anager (`hvm`) is similar to other tools, like `nvm` or `rvm`. It is meant to manage multiple installed versions of Happy and other software packages distributed by CZI via Github.

Each Happy project can create a `.happy/version.lock` JSON file, which specifies which versions of various packages the project requires. `hvm` can automatically install them and manage your path for you via `zsh` shell hooks or explicit CLI commands.

# Compatibility Note

This version of `hvm` was only tested on Mac and Linux. Behavior on Windows is currently unknown. It'll probably work under WSL2, but that has not been tested. Bug reports welcome, as always.

## Quick Start

Common to all shells:
```
mkdir -p $HOME/.czi/bin $HOME/.czi/hooks

# Download `hvm` to somewhere in your $PATH
curl ....... # todo. Host hvm somewhere.

```

Once you have done the above, follow the instructions specific to your shell below to complete setup.

If using `zsh`:
```
# Make sure $HOME/.czi/bin is in your $PATH
echo 'export PATH=$HOME/.czi/bin:$PATH' >> $HOME/.zshrc

# Download the zsh shell hook and install it
curl .... # todo. Need URL.

echo '
source $HOME/.czi/hooks/hvm-hooks.zsh
' >> $HOME/.zshrc

# Now restart your shell session OR execute the below in your current shell to load the hooks, if desired

source $HOME/.czi/hooks/hvm-hooks.zsh
```

## Usage

### help

Help is always available using the `help` command for `hvm` itself and each subcommand.

```
$ hvm help
Manage multiple installed versions of Happy, and facilitate switching between them.

Usage:
  hvm [command]

Available Commands:
  completion    Generate the autocompletion script for the specified shell
  download      Download the specified binary distribution package for Happy
  env           Calculate environment variables for eval() by the calling shell
  help          Help about any command
  install       Install a version of Happy
  list-releases Get list of available releases
  lock          A brief description of your command
  set-default   Symlink the specified version of happy to $HOME/.czi/bin to be used as default

Flags:
  -h, --help     help for hvm
  -t, --toggle   Help message for toggle

Use "hvm [command] --help" for more information about a command.

```

### download

The `download` subcommand simply downloads the tarball for the requested package to the current directory or a path specified by `-p` or `--path`. It is generally meant for use in shell scripts where the default `install` behavior does not fit the use case.

You can specify the `arch` and `os` if you'd like to download binaries that are not for your current architecture or os. The default is to match the current arch/os. 

Valid arch values: amd64 and arm64
Valie os values: linux, darwin, windows (see COMPATIBILITY note above)

```
$ hvm help download

Allow simple download of the tarball/zip file for a specific version of Happy. OS and 
architecture are detected automatically, but can be overridden with the --os and --arch flags.

Usage:
  hvm download [version] [flags]

Flags:
  -a, --arch string   Force architecture (Default: current)
  -h, --help          help for download
  -o, --os string     Force operating system (Default: current)
  -p, --path string   Path to store the downloaded package (default ".")

```

### env

The `env` command outputs shell commands intended to be `eval`ed by your current shell:

```
eval $(hvm env)
```

This currently only manages `$PATH`, but may manage other things in the future. 

It can be run explicitly as above, but is intended to be run via shell hooks each time your directory changes. This allows `hvm` to automatically manage your environment as specified by your Happy configurations.

```
$ hvm help env

Output to STDOUT a list of env vars which should be eval'ed by the calling shell. This is
used to automatically set PATH and other variables via shell hooks.

Usage:
  hvm env [flags]

Flags:
  -h, --help   help for env

```


### install

The `install` command is used to download a package and extract it under `.czi/versions/<package>/<version>`. This is to allow for multiple concurrently-installed versions. This does NOT set your default version. See `set-default` to choose a default version.

```
$ hvm help install
Install a version of Happy to ~/.happy/versions/ and set it as the current version.

Usage:
  hvm install [flags]

Flags:
  -a, --arch string   Force architecture (Default: current)
  -h, --help          help for install
  -o, --os string     Force operating system (Default: current)
  -p, --path string   Path to store the downloaded package (default ".")

```


### list-releases

Get a list of available releases of a package.

```
$ Get list of available releases

Usage:
  hvm list-releases [flags]

Flags:
  -h, --help   help for list-releases

```

### lock

Create/update the locked version of a package in the current happy project. If you are not in a happy project, this will return an error.

```
$ hvm help lock

Lock the current version of happy in the current project. This will create a .happy/version.lock file

Usage:
  hvm lock [flags]

Flags:
  -h, --help   help for lock

```

### set-default

Set a version of a package to be used when:

* Outside of a Happy project
* The current project does not contain a `version.lock` file
* The `version.lock` file does not specify a version for a given package

This works by creating a symbolic link under `$HOME/.czi/bin` to the binary installed under a specific version of the package. 

When inside a Happy project, `hvm` will prepend any specified versions to your `$PATH`, so they will supercede this default. If the project does not specify a version for a given package, or does not have a `.happy/versions.lock` file, there will be no match in your `$PATH`, and the path search will eventually fall back to `$HOME/.czi/bin`.


```
$ hvm help set-default
Create a symbolic link $HOME/.czi/bin/ pointing to the specified version of happy. Assuming
$HOME/.czi/bin is set appropriately in your $PATH, this version will be used by default when running 'happy'
outside of a project, or when a happy version config is not present.

Usage:
  hvm set-default [version] [flags]

Flags:
  -h, --help   help for set-default

```

## Installing the zsh Shell Hook


## Environment Variables

