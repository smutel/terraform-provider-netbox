# netbox_tenancy_tenant Data Source

Get info about tenancy tenant in the netbox provider.

## Example Usage

```hcl
data "netbox_tenancy_tenant" "tenant_test" {
  slug = "TestTenant"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenant.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
