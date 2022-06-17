data "netbox_json_extras_object_changes_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_object_changes_list.test.json)
}
