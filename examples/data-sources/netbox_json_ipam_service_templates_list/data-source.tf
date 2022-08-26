data "netbox_json_ipam_service_templates_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_ipam_service_templates_list.test.json)
}
