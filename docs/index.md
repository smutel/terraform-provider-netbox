# terraform-provider-netbox Provider

Terraform provider for [Netbox.](https://netbox.readthedocs.io/en/stable/)

## Compatibility with Netbox

Version 0.x.x => Netbox 2.8.x

## Example Usage

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

## Argument Reference

* `url` or `NETBOX_URL` environment variable to define the URL and the port (127.0.0.1:8000 by default)
* `token` or `NETBOX_TOKEN` environment variable to define the TOKEN to access the application (empty by default)
* `scheme` or `NETBOX_SCHEME` environment variable to define the SCHEME of the URL (https by default)
