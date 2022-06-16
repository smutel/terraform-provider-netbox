data "netbox_json_ipam_route_targets_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_route_targets_list.test.json)
}
