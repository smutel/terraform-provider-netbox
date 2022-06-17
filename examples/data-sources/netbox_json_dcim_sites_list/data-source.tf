data "netbox_json_dcim_sites_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_sites_list.test.json)
}
