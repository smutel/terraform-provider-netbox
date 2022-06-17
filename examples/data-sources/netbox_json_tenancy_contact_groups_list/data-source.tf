data "netbox_json_tenancy_contact_groups_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_tenancy_contact_groups_list.test.json)
}
