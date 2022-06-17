data "netbox_json_virtualization_cluster_groups_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_virtualization_cluster_groups_list.test.json)
}
