# netbox\_ipam\_ip\_addresses Resource

Manage an ip address within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_ip_addresses" "ip_test" {
  address = "192.168.56.0/24"
  description = "IP created by terraform"
  status = "active"
  
  tag {
    name = "tag1"
    slug = "tag1"
  }
  
  custom_fields = {
    cf_boolean = "true"
    cf_date = "2020-12-25"
    cf_integer = "10"
    cf_selection = "1"
    cf_text = "Some text"
    cf_url = "https://github.com"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``address`` - (Required) The IP address (with mask) used for this object.
* ``custom_fields`` - (Optional) Custom Field Keys and Values for this object
  * For boolean, use the string value "true" or "false"
  * For data, use the string format "YYYY-MM-DD"
  * For integer, use the value between double quote "10"
  * For selection, use the level id
  * For text, use the string value
  * For URL, use the URL as string
* ``description`` - (Optional) The description of this object.
* ``dns_name`` - (Optional) The DNS name of this object.
* ``nat_inside_id`` - (Optional) The ID of the NAT inside of this object.
* ``nat_outside_id`` - (Optional) The ID of the NAT outside of this object.
* ``object_id`` - (Optional) The ID of the object where this resource is attached to.
* ``object_type`` - (Optional) The object type among virtualization.vminterface
or dcim.interface (virtualization.vminterface by default)
* ``primary_ip4`` - (Optional) Set this resource as primary IPv4 (false by default)
* ``role`` - (Optional) The role among loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp of this object.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
