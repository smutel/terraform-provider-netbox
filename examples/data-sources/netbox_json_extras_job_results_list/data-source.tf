data "netbox_json_extras_job_results_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_job_results_list.test.json)
}
