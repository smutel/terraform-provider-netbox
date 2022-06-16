data "netbox_json_ipam_fhrp_group_assignments_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_fhrp_group_assignments_list.test.json)
}
