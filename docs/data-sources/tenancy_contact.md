# netbox\_tenancy\_contact Data Source

Get info about tenancy contact from netbox.

## Example Usage

```hcl
data "netbox_tenancy_contact" "contact_test" {
  name = "John Doe"
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name of the contact.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
