data "netbox_json_extras_config_contexts_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_config_contexts_list.test.json)
}
