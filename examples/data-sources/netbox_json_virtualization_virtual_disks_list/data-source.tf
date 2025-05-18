data "netbox_json_virtualization_virtual_disks_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_virtualization_virtual_disks_list.test.json)
}
