data "netbox_json_ipam_vrfs_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_vrfs_list.test.json)
}
