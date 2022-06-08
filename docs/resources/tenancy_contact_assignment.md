# netbox\_tenancy\_contact\_assignment Resource

Link a contact to another resource within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_contact_assignment" "contact_assignment_01" {
  contact_id = netbox_tenancy_contact.contact.id
  contact_role_id = netbox_tenancy_contact_role.contact_role_02.id
  content_type = netbox_virtualization_vm.vm_test.content_type
  object_id = netbox_virtualization_vm.vm_test.id
  priority = "primary"
}
```

## Argument Reference

The following arguments are supported:
* ``contact_id`` - (Required) ID of the contact to link to a resource.
* ``contact_role_id`` - (Required) The role of the contact for this resource.
* ``content_type`` - (Required) Type of the object where the contact will be linked.
* ``object_id`` - (Required) ID of the object where the contact will be linked.
* ``priority`` - (Required) Priority of this contact among primary, secondary and tertiary (primary by default).

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Contact assignments can be imported by `id` e.g.

```
$ terraform import netbox_tenancy_contact_assignment.contact_assignment_test id
```
