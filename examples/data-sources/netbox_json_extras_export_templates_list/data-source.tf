data "netbox_json_extras_export_templates_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_export_templates_list.test.json)
}
