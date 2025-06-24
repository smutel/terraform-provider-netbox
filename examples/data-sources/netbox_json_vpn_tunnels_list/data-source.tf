data "netbox_json_vpn_tunnels_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_vpn_tunnels_list.test.json)
}
