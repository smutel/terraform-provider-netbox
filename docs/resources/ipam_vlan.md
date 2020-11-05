# netbox\_ipam\_vlan Resource

Manages a vlan within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_vlan" "vlan_test" {
  vlan_id = 100
  name = "TestVlan"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  description = "VLAN created by terraform"
  vlan_group_id = netbox_ipam_vlan_group.vlan_group_test.id
  tenant_id = netbox_tenancy_tenant.tenant_test.id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
  
  tag {
    name = "tag1"
    slug = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``description`` - (Optional) The description of this object.
* ``vlan_group_id`` - (Optional) ID of the group where this object belongs to.
* ``name`` - (Required) The name for this object.
* ``role_id`` - (Optional) The ID of the role attached to this object.
* ``site_id`` - (Optional) ID of the site where this object is created.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vlan_id`` - (Required) The ID of this vlan (vlan tag).

The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
