---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_tenancy_tenant_group Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  
---

# netbox_tenancy_tenant_group (Data Source)



## Example Usage

```terraform
data "netbox_tenancy_tenant_group" "tenancy_tenant_group" {
  slug = "TestTenantGroup"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `slug` (String)

### Read-Only

- `content_type` (String)
- `id` (String) The ID of this resource.


