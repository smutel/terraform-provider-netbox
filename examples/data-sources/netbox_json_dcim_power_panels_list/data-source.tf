data "netbox_json_dcim_power_panels_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_power_panels_list.test.json)
}
