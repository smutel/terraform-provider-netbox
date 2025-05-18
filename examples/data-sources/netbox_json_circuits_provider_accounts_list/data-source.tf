data "netbox_json_circuits_provider_accounts_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_circuits_provider_accounts_list.test.json)
}
