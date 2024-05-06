# Terraform Provider Windows
<!-- Badges -->
[![Build][build badge]][build page]
[![GoReport][goreport badge]][goreport page]
[![Conventional Commits][convention badge]][convention page]

Terraform provider to manage Windows based resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```terraform
terraform {
  required_providers {
    windows = {
      source = "d-strobel/windows"
    }
  }
}

provider "windows" {
  endpoint = "127.0.0.1"

  ssh = {
    username = "vagrant"
    password = "vagrant"
    port     = 1222
  }
}

// Create a new local security group.
resource "windows_local_group" "this" {
    name = "MyNewGroup"
}

// Create a new local user.
resource "windows_local_user" "this" {
    name = "MyNewUser"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

<!-- Badges -->
[goreport badge]: https://goreportcard.com/badge/github.com/d-strobel/terraform-provider-windows
[goreport page]: https://goreportcard.com/report/github.com/d-strobel/terraform-provider-windows

[build badge]: https://github.com/d-strobel/terraform-provider-windows/actions/workflows/build.yml/badge.svg
[build page]: https://github.com/d-strobel/terraform-provider-windows/actions/workflows/build.yml

[convention badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white
[convention page]: https://conventionalcommits.org
