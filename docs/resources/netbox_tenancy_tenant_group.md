# netbox_tenancy_tenant_group Resource

Manages an tenancy tenant group resource within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_tenant_group" "tenant_group_test" {
  name = "TestTenantGroup"
  slug = "TestTenantGroup"
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
