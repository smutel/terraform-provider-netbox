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
  
  custom_field {
    name = "cf_boolean"
    type = "boolean"
    value = "true"
  }

  custom_field {
    name = "cf_date"
    type = "date"
    value = "2020-12-25"
  }

  custom_field {
    name = "cf_text"
    type = "text"
    value = "some text"
  }

  custom_field {
    name = "cf_integer"
    type = "integer"
    value = "10"
  }

  custom_field {
    name = "cf_selection"
    type = "selection"
    value = "1"
  }

  custom_field {
    name = "cf_url"
    type = "url"
    value = "https://github.com"
  }

  custom_field {
    name = "cf_multiple_selection"
    type = "multiple"
    value = "0,1"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``address`` - (Required) The IP address (with mask) used for this object.
* ``description`` - (Optional) The description of this object.
* ``dns_name`` - (Optional) The DNS name of this object.
* ``nat_inside_id`` - (Optional) The ID of the NAT inside of this object.
* ``object_id`` - (Optional) The ID of the object where this resource is attached to.
* ``object_type`` - (Optional) The object type among virtualization.vminterface
or dcim.interface (empty by default)
* ``primary_ip4`` - (Optional) Set this resource as primary IPv4 (false by default)
* ``role`` - (Optional) The role among loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp of this object.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

The ``custom_field`` block (optional) supports:
* ``name`` - (Required) Name of the existing custom resource to associate with this resource.
* ``type`` - (Required) Type of the existing custom resource to associate with this resource (text, integer, boolean, url, selection, multiple).
* ``value`` - (Required) Value of the existing custom resource to associate with this resource.

The ``tag`` block (optional) supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Ip addresses can be imported by `id` e.g.

```
$ terraform import netbox_ipam_ip_addresses.ip_test id
```
