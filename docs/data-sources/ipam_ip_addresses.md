# netbox\_ipam\_ip\_addresses Data Source

Get info about ipam IP addresses from netbox.

## Example Usage

```hcl
data "netbox_ipam_ip_addresses" "ipaddress_test" {
  address = "192.168.56.1/24"
}
```

## Argument Reference

The following arguments are supported:
* ``address`` - (Required) The address (with mask) of the ipam IP address.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
