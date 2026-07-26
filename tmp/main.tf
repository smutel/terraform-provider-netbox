data "netbox_json_dcim_regions_list" "test" {
  limit = 2
  filter {
    name  = "parent"
    value = "test"
  }
}

output "test" {
  value = data.netbox_json_dcim_regions_list.test.json
}
