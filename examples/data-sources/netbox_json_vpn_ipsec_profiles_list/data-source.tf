data "netbox_json_vpn_ipsec_profiles_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_vpn_ipsec_profiles_list.test.json)
}
