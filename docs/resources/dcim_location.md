---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_dcim_location Resource - terraform-provider-netbox"
subcategory: ""
description: |-
  Manage a location (dcim module) within Netbox.
---

# netbox_dcim_location (Resource)

Manage a location (dcim module) within Netbox.

## Example Usage

```terraform
resource "netbox_dcim_location" "location_test" {
  name    = "Test location"
  slug    = "test-location"
  site_id = netbox_dcim_site.test.id

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
    value = 1
  }

  custom_field {
    name = "cf_multi_object"
    type = "multiobject"
    value = jsonencode([
      1,
      2
    ])
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of this location (dcim module).
- `site_id` (Number) The site where this location (dcim module) is.
- `slug` (String) The slug of this site (dcim module).

### Optional

- `custom_field` (Block Set) Existing custom fields to associate to this ressource. (see [below for nested schema](#nestedblock--custom_field))
- `description` (String) Description of this location (dcim module).
- `parent_id` (Number) The ID of the parent for this location (dcim module).
- `status` (String) The status among planned, staging, active, decommissioning or retired (active by default) of this location (dcim module).
- `tag` (Block Set) Existing tag to associate to this resource. (see [below for nested schema](#nestedblock--tag))
- `tenant_id` (Number) The tenant of this location (dcim module).

### Read-Only

- `created` (String) Date when this location (dcim module) was created.
- `depth` (Number) Depth of this location (dcim module).
- `device_count` (Number) Number of devices in this location (dcim module).
- `id` (String) The ID of this resource.
- `last_updated` (String) Date when this location was last updated (dcim module).
- `rack_count` (Number) Number of racks in this location (dcim module).
- `url` (String) The link to this location (dcim module).

<a id="nestedblock--custom_field"></a>
### Nested Schema for `custom_field`

Required:

- `name` (String) Name of the existing custom field.
- `type` (String) Type of the existing custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject, selection (deprecated), multiple(deprecated)).
- `value` (String) Value of the existing custom field.


<a id="nestedblock--tag"></a>
### Nested Schema for `tag`

Required:

- `name` (String) Name of the existing tag.
- `slug` (String) Slug of the existing tag.

## Import

Import is supported using the following syntax:

```shell
# Locations can be imported by id
terraform import netbox_dcim_location.location_test 1
```