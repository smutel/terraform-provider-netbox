# netbox\_virtualization\_vm Resource

Manage a virtual machine resource within Netbox.

## Example Usage

```hcl
resource "netbox_virtualization_vm" "vm_test" {
  name            = "TestVm"
  comments        = "VM created by terraform"
  disk            = 50
  memory          = 16
  cluster_id      = 1

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
* ``cluster_id`` - (Required) ID of the cluster which host this object.
* ``comments`` - (Optional) Comments for this object.
* ``disk`` - (Optional) The size in GB of the disk for this object.
* ``local_context_data`` - (Optional) Local context data for this object.
* ``memory`` - (Optional) The size in MB of the memory of this object.
* ``name`` - (Required) The name for this object.
* ``platform_id`` - (Optional) ID of the platform for this object.
* ``role_id`` - (Optional) ID of the role for this object.
* ``status`` - (Optional) The status among offline, active, planned, staged, failed or decommissioning (active by default).
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vcpus`` - (Optional) The number of VCPUS for this object.
The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.
* ``custom_fields`` - (Optional) Custom Field Keys and Values for this object
** For boolean, use the string value "true" or "false"
** For data, use the string format "YYYY-MM-DD"
** For integer, use the value between double quote "10"
** For selection, use the level id
** For text, use the string value
** For URL, use the URL as string

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
