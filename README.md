# GVM
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/xvrzhao/gvm?label=version)
[![GitHub issues](https://img.shields.io/github/issues/xvrzhao/gvm)](https://github.com/xvrzhao/gvm/issues)
[![GitHub license](https://img.shields.io/github/license/xvrzhao/gvm)](https://github.com/xvrzhao/gvm/blob/master/LICENSE)

GVM is a Go version manager written in Go language, it provides some useful commands like `install`, `list`, `switch` and `remove` to manage local installation of multiple Go versions.

<img src="usage.gif" alt="usage" width="40%" height="40%" />

GVM is just support for Unix-like system yet, and the working mechanism is very simple. 
GVM will create a `gvm` directory in `/usr/local` to host multiple versions of GOROOT, 
and create a symbol link named `go` in `/usr/local` referring to the specific version in `gvm` directory. 
So, you just need to add `/usr/local/go/bin` to `$PATH` environment variable to run go command, 
and use GVM to switch the reference of the symbol link.

Multiple versions of Go installed by GVM can share the same GOPATH compatibly, and this is also advocated by GVM.

## Installation

There are two ways to install GVM.

### Install by Go

If you have installed Go before, just execute the following command:

```
$ go install github.com/xvrzhao/gvm
```

**Note**: 

This will install the GVM binary into your `$GOBIN` (same as `$GOPATH/bin`) directory. If you have added `$GOBIN` to `$PATH`, you can use GVM commands directly. However, some subcommands (like `switch`, `install`, etc.) need to write files in `/usr/local/`, so please make sure you have the appropriate permissions. You can execute the GVM commands with root privilege, like `sudo gvm [command]`.

But sometimes it may prompt `sudo: gvm: command not found`, that is, the root user cannot find GVM binary in his/her `$PATH` directories. Because `sudo` does not use shell login configurations (`/etc/profile`, `$HOME/.bashrc`, etc.) to initialize the `$PATH` environment variable, `$GOBIN` is not the part of `$PATH`. Therefore, when the current user is not `root`, you can use GVM with `sudo $(which gvm) [command]`. Or, thoroughly, install and use GVM under `root` user login.

### Download the binary

Another way is downloading the binary file corresponding to your operating system in the [Releases Page](https://github.com/xvrzhao/gvm/releases).

## Commands

For examples:

```bash
# install and switch to go1.16.3, `--cn` is required for Mainland China.
$ sudo gvm install 1.16.3 --switch --cn 

# list all versions managed by GVM.
$ gvm list

# remove go1.16.3
$ gvm remove 1.16.3

# switch to go1.17
$ gvm switch 1.17
```

For all available commands and flags, see:

```
$ gvm help [subcommand]
```
