data "netbox_json_dcim_modules_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_modules_list.test.json)
}
