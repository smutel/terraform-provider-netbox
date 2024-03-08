data "netbox_ipam_aggregate" "aggregate_test" {
  prefix = "192.168.56.0/24"
  rir_id = 1
}

output "test" {
  value = data.netbox_ipam_aggregate.aggregate_test.id
}
