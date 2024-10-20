# Terraform Provider Windows

<!-- Badges -->
[![Build][build badge]][build page]
[![GoReport][goreport badge]][goreport page]
[![Conventional Commits][convention badge]][convention page]

This Terraform provider enables the management of Windows-based resources within your infrastructure. 

It is built on top of the [gowindows](https://github.com/d-strobel/gowindows) SDK, 
which acts as the underlying interface for interacting with Windows environments.

To introduce new features or enhancements to this provider, the corresponding functionality must first be implemented in the gowindows library. 
Contributions or feature requests should therefore begin with updates to the GoWindows SDK before being integrated into the provider.

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

### Pre-commit

To ensure smooth execution in the pipeline and eliminate potential linting errors,
it's highly advisable to integrate pre-commit hooks. These hooks can be effortlessly
installed to streamline the process and maintain code quality standards.

You can find more details about pre-commit hooks on their official website: [pre-commit](https://pre-commit.com/).

### Conventional Commits

This Terraform provider follows the conventional commit guidelines. For more information, see [conventionalcommits.org](https://www.conventionalcommits.org/).

### Acceptance test

> The acceptance tests are currently not available via Github action.

#### Prerequisites

* [Terraform](https://developer.hashicorp.com/terraform/downloads)
* [OpenTofu](https://opentofu.org/docs/intro/install)
* [Go](https://golang.org/doc/install)
* [Hashicorp Vagrant](https://www.vagrantup.com/)
* [Oracle VirtualBox](https://www.virtualbox.org/)

#### Usage

Boot the Vagrant machines:

```shell
make vagrant-up
```

Run acceptance tests for terraform and opentofu:

```shell
make testacc
```

Destroy the Vagrant machines:

```shell
make vagrant-down
```

###

## Inspirations

* [hashicorp - terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad):<br>

Hashicorp made a great start with the terraform-provider-ad. Currently, it seems that the provider is not actively maintained.<br>
Beyond that, my goal is to split the terraform-provider into a library and a provider and extend its functionality with non Active-Directory systems.

## License
This project is licensed under the [Mozilla Public License Version 2.0](LICENSE).

<!-- Badges -->
[goreport badge]: https://goreportcard.com/badge/github.com/d-strobel/terraform-provider-windows
[goreport page]: https://goreportcard.com/report/github.com/d-strobel/terraform-provider-windows

[build badge]: https://github.com/d-strobel/terraform-provider-windows/actions/workflows/build.yml/badge.svg
[build page]: https://github.com/d-strobel/terraform-provider-windows/actions/workflows/build.yml

[convention badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white
[convention page]: https://conventionalcommits.org
