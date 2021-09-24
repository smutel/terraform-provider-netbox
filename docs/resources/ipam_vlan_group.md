# netbox\_ipam\_vlan\_group Resource

Manage a vlan group within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_vlan_group" "vlan_group_test" {
  name = "TestVlanGroup"
  slug = "TestVlanGroup"
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Vlan groups can be imported by `id` e.g.

```
$ terraform import netbox_ipam_vlan_group.vlan_group_test id
```
