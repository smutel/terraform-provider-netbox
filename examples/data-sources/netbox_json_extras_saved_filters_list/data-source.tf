data "netbox_json_extras_saved_filters_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_saved_filters_list.test.json)
}
