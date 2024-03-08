data "netbox_json_vpn_tunnel_terminations_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_vpn_tunnel_terminations_list.test.json)
}
