data "netbox_json_ipam_asns_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_asns_list.test.json)
}
