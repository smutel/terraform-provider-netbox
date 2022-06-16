data "netbox_json_wireless_wireless_links_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_wireless_wireless_links_list.test.json)
}
