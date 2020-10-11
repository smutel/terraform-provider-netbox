# netbox_ipam_prefix Resource

Manages an ipam prefix resource within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_prefix" "prefix_test" {
  prefix = "192.168.56.0/24"
  vlan_id = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
  tags = ["tag1"]
  status = "active"
}
```

## Argument Reference

The following arguments are supported:
* ``description`` - (Optional) The description of this object.
* ``is_pool`` - (Optional) Define if this object is a pool (false by default).
* ``prefix`` - (Required) The prefix (IP address/mask) used for this object.
* ``role_id`` - (Optional) The ID of the role attached to this object.
* ``site_id`` - (Optional) ID of the site where this object is created
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tags`` - (Optional) Array of tags for this object.
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vlan_id`` - (Optional) ID of the vlan where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

