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

Virtualization vms can be imported by `id` e.g.

```
$ terraform import netbox_virtualization_vm.vm_test id
```
