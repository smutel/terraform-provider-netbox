data "netbox_json_ipam_vlan_groups_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_vlan_groups_list.test.json)
}
