data "netbox_json_core_jobs_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_core_jobs_list.test.json)
}
