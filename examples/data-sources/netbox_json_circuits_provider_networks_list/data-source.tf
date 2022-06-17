data "netbox_json_circuits_provider_networks_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_circuits_provider_networks_list.test.json)
}
