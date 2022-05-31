# netbox\_tenancy\_contact\_role Data Source

Get info about tenancy contact roles from netbox.

## Example Usage

```hcl
data "netbox_tenancy_contact_role" "contact_role_test" {
  slug = "TestContactGroup"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenancy contact roles.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
