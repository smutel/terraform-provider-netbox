# netbox\_tenancy\_contact\_group Data Source

Get info about tenancy contact groups from netbox.

## Example Usage

```hcl
data "netbox_tenancy_contact_group" "contact_group_test" {
  slug = "TestContactGroup"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenancy contact groups.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
