resource "netbox_dcim_platform" "test" {
  name = "test1"
  slug = "test1"
}

resource "netbox_extras_custom_field" "test" {
  name = "test1"
  content_types = [
    "dcim.site",
  ]
  type          = "object"
  object_type   = "dcim.platform"
  description   = "Test custom field"
  group_name    = "testgroup"
  ui_visibility = "hidden"
  label         = "Test Label for CF"
  weight        = 50
  #required      = true
  filter_logic = "disabled"
  #default = jsonencode(
  #  netbox_dcim_platform.test.id
  #)
}
