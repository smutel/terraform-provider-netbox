# netbox\_tenancy\_tenant\_group Data Source

Get info about tenancy tenant groups in the netbox provider.

## Example Usage

```hcl
data "netbox_tenancy_tenant_groups" "tenant_group_test" {
  slug = "TestTenantGroup"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenancy tenant groups.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
