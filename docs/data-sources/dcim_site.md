---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_dcim_site Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  Get info about site (dcim module) from netbox.
---

# netbox_dcim_site (Data Source)

Get info about site (dcim module) from netbox.

## Example Usage

```terraform
data "netbox_dcim_site" "site_test" {
  slug = "TestSite"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `slug` (String) The slug of the site (dcim module).

### Read-Only

- `content_type` (String) The content type of this site (dcim module).
- `id` (String) The ID of this resource.


