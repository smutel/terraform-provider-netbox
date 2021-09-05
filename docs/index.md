# terraform-provider-netbox Provider

Terraform provider for [Netbox.](https://netbox.readthedocs.io/en/stable/)

## Compatibility with Netbox

Version 0.x.y => Netbox 2.8
Version 1.x.y => Netbox 2.9

## Example Usage

```hcl
provider netbox {
  # Environment variable NETBOX_URL
  url = "127.0.0.1:8000"

  # Environment variable NETBOX_BASEPATH
  basepath = "/api"

  # Environment variable NETBOX_TOKEN
  token = "c07a2db4adb8b1e7f75e7c4369964e92f7680512"

  # Environment variable NETBOX_SCHEME
  scheme = "http"

  # Environment variable NETBOX_PRIVATE_KEY_FILE
  private_key_file = "/path/to/private/key"
}
```

## Argument Reference

* `url` or `NETBOX_URL` environment variable to define the URL and the port (127.0.0.1:8000 by default)
* `basepath` or `NETBOX_BASEPATH` environment variable to define the base path (/api)
* `token` or `NETBOX_TOKEN` environment variable to define the TOKEN to access the application (empty by default)
* `scheme` or `NETBOX_SCHEME` environment variable to define the SCHEME of the URL (https by default)
* `insecure` or `NETBOX_INSECURE` environment variable to skip or not the TLS certificat validation (false by default)
* `private_key_file` or `NETBOX_PRIVATE_KEY_FILE` environment variable to add a private key to work with encoded data like secrets (empty by default)
