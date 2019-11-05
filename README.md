# terraform-provider-netbox

[![GoDoc](http://godoc.org/github.com/smutel/terraform-provider-netbox?status.svg)](http://godoc.org/github.com/smutel/terraform-provider-netbox)
[![Build Status](https://github.com/smutel/terraform-provider-netbox/workflows/ci/badge.svg?branch=master)](https://github.com/smutel/terraform-provider-netbox/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/smutel/terraform-provider-netbox)](https://goreportcard.com/report/github.com/smutel/terraform-provider-netbox)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

Terraform provider for [Netbox.](https://netbox.readthedocs.io/en/stable/)

## Requirements

* General developper tools like make, bash, ...
* Go 1.13 minimum (to build the provider)
* Terraform (to use the provider)

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

You can install the provider manually in your global terraform provider folder 
or you can also use the makefile to install the provider in your local provider folder:

```bash
$ make localinstall
==> Creating folder terraform.d/plugins/linux_amd64
mkdir -p ~/.terraform.d/plugins/linux_amd64
==> Installing provider in this folder
cp terraform-provider-netbox ~/.terraform.d/plugins/linux_amd64
```

## Using the provider

The definition of the provider is optional.  
All the paramters could be setup by environment variables.  

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

For further information, check this [documentation](docs/Provider.md)

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
## Known bugs which can impact this provider

N/A
