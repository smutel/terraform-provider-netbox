data "netbox_json_tenancy_tenants_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_tenancy_tenants_list.test.json)
}
