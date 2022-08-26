# go-netbox

[![Lisence](https://img.shields.io/badge/license-ISC-informational?style=flat-square)](https://github.com/smutel/go-netbox/blob/main/LICENSE)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-informational.svg?style=flat-square&logo=git)](https://conventionalcommits.org)
[![Build Status](https://img.shields.io/github/workflow/status/smutel/go-netbox/Main/main?style=flat-square&logo=github-actions)](https://github.com/smutel/go-netbox/actions)

Go library to interact with NetBox IPAM and DCIM service. 

## Compatibility with Netbox

The version for Netbox and go-netbox will the same except for the last digit.

Example:
* go-netbox v1.10.x is working with Netbox v1.10
* go-netbox v1.11.x is working with Netbox v1.11

## Using the client

The `github.com/smutel/go-netbox/netbox` package has some convenience functions for creating clients with the most common
configurations you likely need while connecting to NetBox. `NewNetboxAt` allows you to specify a hostname
(including port, if you need it), and `NewNetboxWithAPIKey` allows you to specify both a hostname:port and API token.

```go
import (
    "github.com/smutel/go-netbox/netbox"
)
...
    c := netbox.NewNetboxAt("your.netbox.host:8000")
    // OR
    c := netbox.NewNetboxWithAPIKey("your.netbox.host:8000", "your_netbox_token")
```

If you specify the API key, you do not need to pass an additional `authInfo` to operations that need authentication, and
can pass `nil`:
```go
    c.Dcim.DcimDeviceTypesCreate(createRequest, nil)
```

If you connect to netbox via HTTPS you have to create an HTTPS configured transport:
```go
package main

import (
	"os"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/dcim"

	log "github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("NETBOX_TOKEN")
	if token == "" {
		log.Fatalf("Please provide netbox API token via env var NETBOX_TOKEN")
	}

	netboxHost := os.Getenv("NETBOX_HOST")
	if netboxHost == "" {
		log.Fatalf("Please provide netbox host via env var NETBOX_HOST")
	}

	transport := httptransport.New(netboxHost, client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+token)

	c := client.New(transport, nil)

	req := dcim.NewDcimSitesListParams()
	res, err := c.Dcim.DcimSitesList(req, nil)
	if err != nil {
		log.Fatalf("Cannot get sites list: %v", err)
	}
	log.Infof("res: %v", res)
}
```

## How to contribute to this project

* To contribute to this project, please follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0-beta.2/) rules.
* Most of the code of this project will be generated using the swagger spec of Netbox and the [go-swagger](https://github.com/go-swagger/go-swagger) program.
* You can change the behavior of the generated library by pushing patchs in the `patchs` folder.
* The best is to see if the bug is due to a wrong swagger definition and to report this bug to the [Netbox](https://github.com/netbox-community/netbox) project.
* If the bug is due to the go-swagger program the best is to create a bug here [go-swagger](https://github.com/go-swagger/go-swagger).

## How to test your work locally

### Requirements

* docker
* docker-compose
* [swagger](https://github.com/go-swagger/go-swagger) installed somewhere (/usr/local/bin)
* [netbox-docker](https://github.com/netbox-community/netbox-docker.git) project installed somewhere

### Installing the go-netbox

```sh
$ # Install jsonlint
$ apt-get install python3-demjson
$ mkdir -p ~/go/src/github.com/smutel
$ cd ~/go/src/github.com/smutel
$ git clone git@github.com:smutel/go-netbox.git
$ export GITHUB_WORKSPACE=~/go/src
```

### Regenerating the library

```sh
$ cd ~/go/src/github.com/smutel/go-netbox/utils
$ ./netbox_generate_client
```
