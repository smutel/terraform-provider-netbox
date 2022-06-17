data "netbox_json_users_tokens_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_users_tokens_list.test.json)
}
