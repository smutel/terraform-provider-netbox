data "netbox_json_dcim_virtual_device_contexts_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_virtual_device_contexts_list.test.json)
}
