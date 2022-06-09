# netbox\_json\_wireless\_wireless\_lans\_list Data Source

Get json output from the wireless_wireless_lans_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_wireless_wireless_lans_list" "test" {
  limit = 0
}
output "example" {
  value = jsondecode(data.netbox_json_wireless_wireless_lans_list.test.json)
}
```

## Argument Reference

* ``limit`` (Optional). The max number of returned results. If 0 is specified, all records will be returned.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

