data "netbox_json_extras_custom_links_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_custom_links_list.test.json)
}
