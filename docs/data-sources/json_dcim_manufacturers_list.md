# netbox\_json\_dcim\_manufacturers\_list Data Source

Get json output from the dcim_manufacturers_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_dcim_manufacturers_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_dcim_manufacturers_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

