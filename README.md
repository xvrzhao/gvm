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

```
$ go get github.com/xvrzhao/gvm
```

This will install the `gvm` bin command into your `$GOPATH/bin` directory.

## Commands

For all available commands, see:

```
$ gvm help
```
