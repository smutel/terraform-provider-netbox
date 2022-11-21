resource "netbox_extras_custom_field" "cf_text" {
  name = "cf_text"
  label = "CF Text"
  type = "text"
  content_types = [
    "dcim.site",
  ]

  description = "Test text field"
  weight = 50
  required = true
  filter_logic = "exact"
  default = "Default value"
  validation_regex = "^.*$"
}

resource "netbox_extras_custom_field" "cf_boolean" {
  name = "cf_boolean"
  label = "CF Boolean"
  type = "boolean"
  content_types = [
    "dcim.site",
  ]

  description = "Test boolean field"
  weight = 50
  filter_logic = "exact"
  # Fixed in Netbox 3.3
  # required = true
  # default = jsonencode(false)
}

resource "netbox_extras_custom_field" "cf_integer" {
  name = "cf_integer"
  label = "CF Integer"
  type = "integer"
  content_types = [
    "dcim.site",
  ]

  description = "Test integer field"
  weight = 50
  filter_logic = "loose"
  # Fixed in Netbox 3.3
  # required = true
  # default = jsonencode(5)
}

resource "netbox_extras_custom_field" "cf_select" {
  name = "cf_select"
  label = "CF Select"
  type = "select"
  content_types = [
    "dcim.site",
  ]

  choices = [
    "choice 1",
    "choice 2",
  ]

  description = "Test select field"
  weight = 50
  filter_logic = "loose"
  required = true
  default = "choice 1"
}

resource "netbox_extras_custom_field" "cf_multi_select" {
  name = "cf_mulit_select"
  label = "CF Multiselect"
  type = "multiselect"
  content_types = [
    "dcim.site",
  ]

  choices = [
    "choice 1",
    "choice 2",
  ]

  description = "Test multiselect field"
  weight = 50
  filter_logic = "loose"
  # Fixed in netbox 3.3
  # required = true
  # default = jsonencode([
  #   "choice 1",
  # ])
}

resource "netbox_extras_custom_field" "cf_object" {
  name = "cf_object"
  label = "CF object"
  type = "object"
  content_types = [
    "dcim.site",
  ]

  object_type = "dcim.platform"
  description = "Test object field"
  weight = 50
  filter_logic = "loose"
  # Fixed in netbox 3.3
  # required = true
  #default = jsonencode(netbox_dcim_platform.test.id)
}

resource "netbox_extras_custom_field" "cf_multi_object" {
  name = "cf_mulit_object"
  label = "CF Multiobject"
  type = "multiobject"
  content_types = [
    "dcim.site",
  ]

  object_type = "dcim.platform"
  description = "Test multiobject field"
  weight = 50
  filter_logic = "loose"
  # Fixed in netbox 3.3
  # required = true
  # default = jsonencode([
  #   netbox_dcim_platform.test.id,
  # ])
}

resource "netbox_extras_custom_field" "cf_json" {
  name = "cf_json"
  label = "CF Json"
  type = "json"
  content_types = [
    "dcim.site",
  ]

  description = "Test multiobject field"
  weight = 50
  filter_logic = "loose"
  # Fixed in netbox 3.3
  # required = true
  #default = jsonencode({
  #	bool = false
  #	number = 1.5
  #	dict = {
  #		text = "Some text"}
  #	}
  #})
}
