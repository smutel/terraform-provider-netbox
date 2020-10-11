# netbox_dcim_site Data Source

Get info about dcim site in the netbox provider.

## Example Usage

```hcl
data "netbox_dcim_site" "site_test" {
  slug = "TestSite"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the dcim site.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
