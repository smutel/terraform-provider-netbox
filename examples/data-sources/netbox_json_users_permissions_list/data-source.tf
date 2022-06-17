data "netbox_json_users_permissions_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_users_permissions_list.test.json)
}
