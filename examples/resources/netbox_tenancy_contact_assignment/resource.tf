resource "netbox_tenancy_contact_assignment" "contact_assignment_01" {
  contact_id = netbox_tenancy_contact.contact.id
  contact_role_id = netbox_tenancy_contact_role.contact_role_02.id
  content_type = netbox_virtualization_vm.vm_test.content_type
  object_id = netbox_virtualization_vm.vm_test.id
  priority = "primary"
}
