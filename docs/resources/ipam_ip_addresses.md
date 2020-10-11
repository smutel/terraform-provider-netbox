# netbox\_ipam\_ip\_addresses Resource

Manages an ipam ip addresses resource within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_ip_addresses" "ip_test" {
  address = "192.168.56.0/24"
  description = "IP created by terraform"
  tags = ["tag1"]
  status = "active"
}
```

## Argument Reference

The following arguments are supported:
* ``address`` - (Required) The IP address (with mask) used for this object.
* ``description`` - (Optional) The description of this object.
* ``dns_name`` - (Optional) The DNS name of this object.
* ``interface_id`` - (Optional) The ID of the interface where this object is attached to.
* ``nat_inside_id`` - (Optional) The ID of the NAT inside of this object.
* ``nat_outside_id`` - (Optional) The ID of the NAT outside of this object.
* ``role`` - (Optional) The role among loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp of this object.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tags`` - (Optional) Array of tags for this object.
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

