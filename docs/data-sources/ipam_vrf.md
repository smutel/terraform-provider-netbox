---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_ipam_vrf Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  Get info about vrf (ipam module) from netbox.
---

# netbox_ipam_vrf (Data Source)

Get info about vrf (ipam module) from netbox.

## Example Usage

```terraform
data "netbox_ipam_vrf" "vrf_test" {
  vrf_id = 15
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `vrf_id` (Number) The ID of the vrf (ipam module).

### Read-Only

- `content_type` (String) The content type of this vrf (ipam module).
- `id` (String) The ID of this resource.


