data "netbox_json_circuits_circuit_terminations_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_circuits_circuit_terminations_list.test.json)
}
