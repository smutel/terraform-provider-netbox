data "netbox_json_extras_journal_entries_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_journal_entries_list.test.json)
}
