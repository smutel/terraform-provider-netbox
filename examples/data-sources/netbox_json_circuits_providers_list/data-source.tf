data "netbox_json_circuits_providers_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_circuits_providers_list.test.json)
}
