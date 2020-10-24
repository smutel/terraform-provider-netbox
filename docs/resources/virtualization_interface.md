# netbox\_virtualization\_interface Resource

Manage an interface resource within Netbox.

## Example Usage

```hcl
resource "netbox_virtualization_interface" "interface_test" {
  name = "default"
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  mac_address = "AA:AA:AA:AA:AA:AA"
  mtu = 1500
  description = "Interface de test"
}
```

## Argument Reference

The following arguments are supported:
* ``description`` - (Optional) Description for this object.
* ``enabled`` - (Optional) true or false (true by default).
* ``mac_address`` - (Optional) Mac address for this object.
* ``mode`` - (Optional) The mode among access, tagged, tagged-all.
* ``mtu`` - (Optional) The MTU between 1 and 65536 for this object.
* ``name`` - (Required) The name for this object.
* ``tagged_vlans`` - (Optional) List of vlan id tagged for this interface
* ``untagged_vlans`` - (Optional) Vlan id untagged for this interface
* ``virtualmachine_id`` - (Required) ID of the virtual machine where this object
is attached
The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
