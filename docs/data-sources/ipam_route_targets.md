---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_ipam_route_targets Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  Get info about vrf (ipam module) from netbox.
---

# netbox_ipam_route_targets (Data Source)

Get info about vrf (ipam module) from netbox.

## Example Usage

```terraform
data "netbox_ipam_route_targets" "rt_test" {
  name = "rt-test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of this Route Targets (ipam module).

### Read-Only

- `content_type` (String) The content type of this vrf (ipam module).
- `id` (String) The ID of this resource.


