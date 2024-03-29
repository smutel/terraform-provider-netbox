---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_ipam_vrf Resource - terraform-provider-netbox"
subcategory: ""
description: |-
  Manage a vrf (ipam module) within Netbox.
---

# netbox_ipam_vrf (Resource)

Manage a vrf (ipam module) within Netbox.

## Example Usage

```terraform
resource "netbox_ipam_vrf" "vrf_test" {
  name = "Test VRF"
  enforce_unique = false
  export_targets = [ netbox_ipam_route_targets.rt_export_test.id ]
  import_targets = [ netbox_ipam_route_targets.rt_import_test.id ]
  rd = "test-vrf"
  description = "Test VRF"
	comments = <<-EOT
	Test Vrf
	EOT

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

- `name` (String) The name of this VRF (ipam module).

### Optional

- `comments` (String) Comments for this VRF (ipam module).
- `custom_field` (Block Set) Existing custom fields to associate to this ressource. (see [below for nested schema](#nestedblock--custom_field))
- `description` (String) The description of this VRF (ipam module).
- `enforce_unique` (Boolean) Prevent duplicate prefixes/IP addresses within this VRF (ipam module)
- `export_targets` (List of Number) Array of ID of exported vrf targets attached to this VRF (ipam module).
- `import_targets` (List of Number) Array of ID of imported vrf targets attached to this VRF (ipam module).
- `rd` (String) The Route Distinguisher (RFC 4364) of this VRF (ipam module).
- `tag` (Block Set) Existing tag to associate to this resource. (see [below for nested schema](#nestedblock--tag))
- `tenant_id` (Number) ID of the tenant where this VRF (ipam module) is attached.

### Read-Only

- `content_type` (String) The content type of this VRF (ipam module).
- `created` (String) Date when this VRF was created.
- `id` (String) The ID of this resource.
- `last_updated` (String) Date when this VRF was created.
- `url` (String) The link to this VRF (ipam module).

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
# VRFs can be imported by id
terraform import netbox_ipam_vrf.vrf_test 1
```
