data "netbox_json_vpn_ike_policies_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_vpn_ike_policies_list.test.json)
}
