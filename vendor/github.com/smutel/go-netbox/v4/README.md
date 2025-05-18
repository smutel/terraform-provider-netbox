# go-netbox

[![Lisence](https://img.shields.io/badge/license-ISC-informational?style=flat-square)](https://github.com/smutel/go-netbox/blob/main/LICENSE)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-informational.svg?style=flat-square&logo=git)](https://conventionalcommits.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/smutel/go-netbox/main.yml?branch=main&style=flat-square)](https://github.com/smutel/go-netbox/actions)

Go library to interact with NetBox IPAM and DCIM service. 

## Compatibility with Netbox

The version for Netbox and go-netbox will the same except for the last digit.

Example:
* go-netbox v2.11.x is working with Netbox v2.11
* go-netbox v3.4.x is working with Netbox v3.4
* go-netbox v4.0.x is working with Netbox v4.0

## Using the client

As an example of usage of this library, you can check this [project](https://github.com/smutel/terraform-provider-netbox)

## How to contribute to this project

* To contribute to this project, please follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0-beta.2/) rules.
* Most of the code of this project will be generated using the openapi spec of Netbox and the [openapi-generator-cli](https://github.com/OpenAPITools/openapi-generator-cli) program.
* You can change the behavior of the generated library by updating python script `utils/fix-spec.py` or by pushing patchs in the `patches` folder.
* The best is to see if the bug is due to a wrong openapi definition and to report this bug to the [Netbox](https://github.com/netbox-community/netbox) project.
* If the bug is due to the `openapi-generator-cli` program the best is to create a bug here [openapi-generator-cli](https://github.com/OpenAPITools/openapi-generator-cli).
* If the bug is due to the `openapi-generator` program the best is to create a bug here [openapi-generator](https://github.com/OpenAPITools/openapi-generator).

## How to test your work locally

### Requirements

* docker
* [netbox-docker](https://github.com/netbox-community/netbox-docker.git) project installed somewhere

### Installing the go-netbox

```sh
$ mkdir -p ~/go/src/github.com/smutel
$ cd ~/go/src/github.com/smutel
$ git clone git@github.com:smutel/go-netbox.git
```

### Regenerating the library

```sh
$ cd ~/go/src/github.com/smutel/go-netbox/utils
$ export GITHUB_WORKSPACE=~/go/src
$ ./netbox_generate_client
```
