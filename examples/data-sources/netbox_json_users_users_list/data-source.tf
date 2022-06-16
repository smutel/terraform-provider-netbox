data "netbox_json_users_users_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_users_users_list.test.json)
}
