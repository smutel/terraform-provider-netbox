data "netbox_json_dcim_devices_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_devices_list.test.json)
}
