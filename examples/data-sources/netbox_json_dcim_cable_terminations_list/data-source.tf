data "netbox_json_dcim_cable_terminations_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_cable_terminations_list.test.json)
}
