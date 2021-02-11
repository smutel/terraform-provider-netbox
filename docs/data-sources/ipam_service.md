# netbox\_ipam\_service Data Source

Get info about ipam service from netbox.

## Example Usage

```hcl
data "netbox_ipam_service" "service_test" {
  device_id = 5
  name      = "Mail"
  port      = 25
  protocol  = "tcp"
}
```

## Argument Reference

The following arguments are supported:
* ``device_id`` - (Optional) The ID of the device linked to this object.
* ``name`` - (Required) The name of this object.
* ``port`` - (Required) The port of this object.
* ``protocol`` - (Required) The protocol of this service (tcp or udp).
* ``virtualmachine_id`` - (Optional) The ID of the vm linked to this object.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
