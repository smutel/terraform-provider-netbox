data "netbox_json_extras_custom_field_choice_sets_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_custom_field_choice_sets_list.test.json)
}
