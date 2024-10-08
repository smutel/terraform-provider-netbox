data "netbox_json_extras_scripts_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_scripts_list.test.json)
}
