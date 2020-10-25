# terraform-provider-netbox

[![Lisence](https://img.shields.io/badge/license-ISC-informational?style=flat-square)](https://github.com/smutel/terraform-provider-netbox/blob/master/LICENSE)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-informational.svg?style=flat-square&logo=git)](https://conventionalcommits.org)
[![Build Status](https://img.shields.io/github/workflow/status/smutel/terraform-provider-netbox/checks/master?style=flat-square&logo=github-actions)](https://github.com/smutel/terraform-provider-netbox/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/smutel/terraform-provider-netbox?style=flat-square)](https://goreportcard.com/report/github.com/smutel/terraform-provider-netbox)

Terraform provider for [Netbox.](https://netbox.readthedocs.io/en/stable/)

## Requirements

* General developper tools like make, bash, ...
* Go 1.14 minimum (to build the provider)
* Terraform (to use the provider)

## Compatibility with Netbox

Version 0.x.y => Netbox 2.8  
Version 1.x.y => Netbox 2.9  

## Building the provider

Clone repository to: ``$GOPATH/src/github.com/smutel/terraform-provider-netbox``

```bash
$ mkdir -p $GOPATH/src/github.com/smutel
$ cd $GOPATH/src/github.com/smutel
$ git clone git@github.com:smutel/terraform-provider-netbox.git
```

Enter the provider directory and build the provider

```bash
$ cd $GOPATH/src/github.com/smutel/terraform-provider-netbox
$ make build
```

## Installing the provider

---
**NOTE**

Before changing the version of the provider, please remove the temporary folder `.terraform` and `~/.terraform.d`.

---

### Automatic installation from Terraform 0.13

```hcl
terraform {
  required_providers {
    netbox = {
      source = "smutel/netbox"
      version = "1.1.0"
    }
  }
}
```

### Manual installation

You can install the provider manually in your global terraform provider folder.

```bash
$ export NETBOX_PROVIDER_VERSION=1.1.0
$ mkdir -p ~/.terraform.d/plugins/registry.terraform.io/smutel/netbox/${NETBOX_PROVIDER_VERSION}/linux_amd64
$ cp terraform-provider-netbox_v${NETBOX_PROVIDER_VERSION} ~/.terraform.d/plugins/registry.terraform.io/smutel/netbox/${NETBOX_PROVIDER_VERSION}/linux_amd64/terraform-provider-netbox_v${NETBOX_PROVIDER_VERSION}
```

### Manual installation (to test the compiled version = version 1.1.0)

```bash
$ make localinstall
==> Creating folder ~/.terraform.d/plugins/registry.terraform.io/smutel/netbox/1.1.0/linux_amd64
==> Installing provider in this folder
```

## Using the provider

The definition of the provider is optional.  
All the parameters could be setup by environment variables.  

```hcl
provider netbox {
  # Environment variable NETBOX_URL
  url = "127.0.0.1:8000"
  
  # Environment variable NETBOX_TOKEN
  token = "c07a2db4adb8b1e7f75e7c4369964e92f7680512"
  
  # Environment variable NETBOX_SCHEME
  scheme = "http"
}
```

For further information, check this [documentation](https://registry.terraform.io/providers/smutel/netbox/latest/docs)

## Contributing to this project

To contribute to this project, please follow the [conventional
commits](https://www.conventionalcommits.org/en/v1.0.0-beta.2/) rules.

## Examples

You can find some examples in the examples folder.  
Each example can be executed directly with command terraform init & terraform apply.  
You can set different environment variables for your test:
* NETBOX_URL to define the URL and the port (127.0.0.1:8000 by default)
* NETBOX_TOKEN to define the TOKEN to access the application (empty by default)
* NETBOX_SCHEME to define the SCHEME of the URL (https by default)

```bash
$ export NETBOX_URL="127.0.0.1:8000"
$ export NETBOX_TOKEN="c07a2db4adb8b1e7f75e7c4369964e92f7680512"
$ export NETBOX_SCHEME="http"
$ cd examples
$ terraform init & terraform apply
```
## Known bugs in external project which can impact this provider

* [Issue 85 in project go-netbox](https://github.com/netbox-community/go-netbox/issues/85)
* [Issue 54 in project go-netbox](https://github.com/netbox-community/go-netbox/issues/54)
