data "netbox_json_virtualization_clusters_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_virtualization_clusters_list.test.json)
}
