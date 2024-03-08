data "netbox_json_ipam_asn_ranges_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_asn_ranges_list.test.json)
}
