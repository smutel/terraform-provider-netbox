data "netbox_json_dcim_virtual_chassis_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_virtual_chassis_list.test.json)
}
