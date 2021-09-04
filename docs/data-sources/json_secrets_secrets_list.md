# netbox\_json\_secrets\_secrets\_list Data Source

Get json output from the secrets_secrets_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_secrets_secrets_list" "test" {
  filter {
    name = ""
    value = ""
  }
  limit = 0
}
output "example" {
  value = jsondecode(data.netbox_json_secrets_secrets_list.test.json)
}
```

## Argument Reference

The following arguments are supported:

* ``filter`` (Optional). This block set should include "name" (String) and "value" (String).
  Supported options:
  - ``device_name`` The name of the device.
  - ``tag`` The tag of a secret.
* ``limit`` (Optional). The max number of returned results. If 0 is specified, all secrets will be returned.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

