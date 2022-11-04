resource "netbox_ipam_service" "service_test" {
  name              = "SMTP"
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  ip_addresses_id   = [netbox_ipam_ip_addresses.ip_test.id]
  ports             = ["22"]
  protocol          = "tcp"
  description       = "Service created by terraform"

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_field {
    name = "cf_boolean"
    type = "boolean"
    value = "true"
  }

  custom_field {
    name = "cf_date"
    type = "date"
    value = "2020-12-25"
  }

  custom_field {
    name = "cf_text"
    type = "text"
    value = "some text"
  }

  custom_field {
    name = "cf_integer"
    type = "integer"
    value = "10"
  }

  custom_field {
    name = "cf_selection"
    type = "select"
    value = "1"
  }

  custom_field {
    name = "cf_url"
    type = "url"
    value = "https://github.com"
  }

  custom_field {
    name = "cf_multi_selection"
    type = "multiselect"
    value = jsonencode([
      "0",
      "1"
    ])
  }

  custom_field {
    name = "cf_json"
    type = "json"
    value = jsonencode({
      stringvalue = "string"
      boolvalue = false
      dictionary = {
        numbervalue = 5
      }
    })
  }

  custom_field {
    name = "cf_object"
    type = "object"
    value = data.netbox_dcim_platform.platform_test.id
  }

  custom_field {
    name = "cf_multi_object"
    type = "multiobject"
    value = jsonencode([
      data.netbox_dcim_platform.platform_test.id,
      data.netbox_dcim_platform.platform_test2.id
    ])
  }
}
