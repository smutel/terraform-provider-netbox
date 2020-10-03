# netbox\_tenancy\_tenant Resource

Manages an tenancy tenant resource within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "TestTenant"
  slug            = "TestTenant"
  description     = "Tenant created by terraform"
  comments        = "Some test comments"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id
  
  tag {
    name = "tag1"
    slug = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``comments`` - (Optional) Comments for this object.
* ``description`` - (Optional) The description for this object.
* ``tenant_group_id`` - (Optional) ID of the group where this object is located.
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.
The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
