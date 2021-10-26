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
  
  custom_field {
    name = "cf_boolean"
    type = "boolean"
    value = "true"
  }

  custom_field {
    name = "cf_date"
    type = "date"
    value = "2020-12-25"
  }

  custom_field {
    name = "cf_text"
    type = "text"
    value = "some text"
  }

  custom_field {
    name = "cf_integer"
    type = "integer"
    value = "10"
  }

  custom_field {
    name = "cf_selection"
    type = "selection"
    value = "1"
  }

  custom_field {
    name = "cf_url"
    type = "url"
    value = "https://github.com"
  }

  custom_field {
    name = "cf_multiple_selection"
    type = "multiple"
    value = "0,1"
  }
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

The ``custom_field`` block (optional) supports:
* ``name`` - (Required) Name of the existing custom resource to associate with this resource.
* ``type`` - (Required) Type of the existing custom resource to associate with this resource (text, integer, boolean, url, selection, multiple).
* ``value`` - (Required) Value of the existing custom resource to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Virtualization interfaces can be imported by `id` e.g.

```
$ terraform import netbox_virtualization_interface.interface_test id
```
