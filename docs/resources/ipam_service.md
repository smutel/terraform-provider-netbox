# netbox\_ipam\_service Resource

Manage a service within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_service" "service_test" {
  name              = "SMTP"
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  ip_addresses_id   = [netbox_ipam_ip_addresses.ip_test.id]
  port              = "22"
  protocol          = "tcp"
  description       = "Service created by terraform"

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
  }
}
```

## Argument Reference

The following arguments are supported:
* ``custom_fields`` - (Optional) Custom Field Keys and Values for this object
  * For boolean, use the string value "true" or "false"
  * For data, use the string format "YYYY-MM-DD"
  * For integer, use the value between double quote "10"
  * For selection, use the level id
  * For text, use the string value
  * For URL, use the URL as string
* ``description`` - (Optional) The description of this object.
* ``device_id`` - (Optional) The ID of the device linked to this object.
* ``ip_addresses_id`` - (Optional) Array of ID of the IP addresses attached to this object.
* ``name`` - (Required) The name for this object.
* ``port`` - (Optional) The port of this object.
* ``protocol`` - (Required) The protocol of this object (tcp or udp).
* ``virtualmachine_id`` - (Optional) The ID of the vm linked to this object.

The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Import

Services can be imported by `id` e.g.

```
$ terraform import netbox_ipam_service.service_test id
```
