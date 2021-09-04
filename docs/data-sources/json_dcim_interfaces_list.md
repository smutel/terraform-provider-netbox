# netbox\_json\_dcim\_interfaces\_list Data Source

Get json output from the dcim_interfaces_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_dcim_interfaces_list" "test" {
  devicename = ""
  limit = 0
}
output "example" {
  value = jsondecode(data.netbox_json_dcim_interfaces_list.test.json)
}
```

## Argument Reference

The following arguments are supported:

* ``devicename`` (Optional). The name of the interface device.
* ``limit`` (Optional). The max number of returned results. If 0 is specified, all interfaces will be returned.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

