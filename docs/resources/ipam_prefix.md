# netbox\_ipam\_prefix Resource

Manage a prefix within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_prefix" "prefix_test" {
  prefix = "192.168.56.0/24"
  vlan_id = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
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
* ``description`` - (Optional) The description of this object.
* ``is_pool`` - (Optional) Define if this object is a pool (false by default).
* ``prefix`` - (Required) The prefix (IP address/mask) used for this object.
* ``role_id`` - (Optional) The ID of the role attached to this object.
* ``site_id`` - (Optional) ID of the site where this object is created
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vlan_id`` - (Optional) ID of the vlan where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

The ``custom_field`` block (optional) supports:
* ``name`` - (Required) Name of the existing custom resource to associate with this resource.
* ``type`` - (Required) Type of the existing custom resource to associate with this resource (text, integer, boolean, url, selection, multiple).
* ``value`` - (Required) Value of the existing custom resource to associate with this resource.

The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.

## Import

Prefixes can be imported by `id` e.g.

```
$ terraform import netbox_ipam_prefix.prefix_test id
```
