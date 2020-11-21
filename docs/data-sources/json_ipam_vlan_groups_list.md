# netbox\_json\_ipam\_vlan\_groups\_list Data Source

Get json output from the ipam_vlan_groups_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_ipam_vlan_groups_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_ipam_vlan_groups_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

