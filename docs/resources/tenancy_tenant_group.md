# netbox\_tenancy\_tenant\_group Resource

Manage a tenant group within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_tenant_group" "tenant_group_test" {
  name = "TestTenantGroup"
  slug = "TestTenantGroup"
  
  tag {
    name = "tag1"
    slug = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.

The ``tag`` block (optional) supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Tenant groups can be imported by `id` e.g.

```
$ terraform import netbox_tenancy_tenant_group.tenant_group_test id
```
