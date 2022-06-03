# netbox\_ipam\_vlan Data Source

Get info about ipam vlan from netbox.

## Example Usage

```hcl
data "netbox_ipam_vlan" "vlan_test" {
  vlan_id = 15
  vlan_group_id = 16
}
```

## Argument Reference

The following arguments are supported:
* ``vlan_id`` - (Required) The ID of the ipam vlan.
* ``vlan_group_id`` - (Optional) The id of the vlan group linked to this vlan.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
