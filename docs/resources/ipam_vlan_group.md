# netbox\_ipam\_vlan\_group Resource

Manages an ipam vlan group resource within Netbox.


## Example Usage

```hcl
resource "netbox_ipam_vlan_group" "vlan_group_test" {
  name = "TestVlanGroup"
  slug = "TestVlanGroup"
  site_id = data.netbox_dcim_sites.site_test.id
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``site_id`` - (Optional) ID of the site where this object is created.
* ``slug`` - (Required) The slug for this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
