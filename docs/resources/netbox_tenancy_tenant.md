# netbox_tenancy_tenant Resource

Manages an tenancy tenant resource within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "TestTenant"
  slug            = "TestTenant"
  description     = "Tenant created by terraform"
  comments        = "Some test comments"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id
  tags            = ["tag1"]
}
```

## Argument Reference

The following arguments are supported:
* ``comments`` - (Optional) Comments for this object.
* ``description`` - (Optional) The description for this object.
* ``tenant_group_id`` - (Optional) ID of the group where this object is located.
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.
* ``tags`` - (Optional) Array of tags for this tenant.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
