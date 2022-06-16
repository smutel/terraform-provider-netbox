data "netbox_json_dcim_platforms_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_platforms_list.test.json)
}
