---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_dcim_region Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  Get info about region (dcim module) from netbox.
---

# netbox_dcim_region (Data Source)

Get info about region (dcim module) from netbox.

## Example Usage

```terraform
data "netbox_dcim_region" "region_test" {
  slug = "TestRegion"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `slug` (String) The slug of the region (dcim module).

### Read-Only

- `content_type` (String) The content type of this region (dcim module).
- `id` (String) The ID of this resource.


