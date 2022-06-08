# netbox\_tenancy\_tenant Data Source

Get info about tenancy tenant from netbox.

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
* ``content_type`` - The content type of this object.
