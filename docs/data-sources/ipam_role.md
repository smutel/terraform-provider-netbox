# netbox\_ipam\_role Data Source

Get info about ipam role from netbox.

## Example Usage

```hcl
data "netbox_ipam_role" "role_test" {
  slug = "TestRole"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the ipam role.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
