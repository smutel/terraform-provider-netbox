resource "netbox_virtualization_vm_primary_ip" "name" {
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  primary_ip4_id = netbox_ipam_ip_addresses.ip_test.id
  primary_ip6_id = netbox_ipam_ip_addresses.ip6_test.id
}
