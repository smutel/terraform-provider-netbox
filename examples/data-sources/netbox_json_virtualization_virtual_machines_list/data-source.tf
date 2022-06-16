data "netbox_json_virtualization_virtual_machines_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_virtualization_virtual_machines_list.test.json)
}
