data "netbox_json_ipam_aggregates_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_aggregates_list.test.json)
}
