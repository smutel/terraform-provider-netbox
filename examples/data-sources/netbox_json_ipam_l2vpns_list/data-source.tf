data "netbox_json_ipam_l2vpns_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_l2vpns_list.test.json)
}
