# netbox\_dcim\_platform Data Source

Get info about dcim platform from netbox.

## Example Usage

```hcl
data "netbox_dcim_platform" "platform_test" {
  slug = "TestPlatform"
}
```

## Argument Reference

The following arguments are supported:
* ``slug`` - (Required) The slug of the dcim platform.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
