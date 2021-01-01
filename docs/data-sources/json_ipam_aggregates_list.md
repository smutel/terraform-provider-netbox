# netbox\_json\_ipam\_aggregates\_list Data Source

Get json output from the ipam_aggregates_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_ipam_aggregates_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_ipam_aggregates_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

