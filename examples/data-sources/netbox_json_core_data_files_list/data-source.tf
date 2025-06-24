data "netbox_json_core_data_files_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_core_data_files_list.test.json)
}
