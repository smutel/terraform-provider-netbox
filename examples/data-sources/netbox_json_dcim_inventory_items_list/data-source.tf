data "netbox_json_dcim_inventory_items_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_dcim_inventory_items_list.test.json)
}
