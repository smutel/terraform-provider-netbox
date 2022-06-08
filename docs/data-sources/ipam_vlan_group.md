# netbox\_ipam\_vlan\_group Data Source

Get info about ipam vlan group from netbox.

## Example Usage

```hcl
data "netbox_ipam_vlan_group" "vlan_group_test" {
  slug = "TestVlanGroup"
  site_id = 15
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the ipam vlan group.
* ``site_id`` - (Optional) The site_id of the ipam vlan groups.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
