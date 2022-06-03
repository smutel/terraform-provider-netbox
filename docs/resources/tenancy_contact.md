# netbox\_tenancy\_contact Resource

Manage a contact within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_contact" "contact_test" {
  name = "John Doe"
  title = "Someone in the world"
  phone = "+330123456789"
  email = "john.doe@unknown.com"
  address = "Somewhere in the world"
  comments = "Good contact"
  contact_group_id = netbox_tenancy_contact_group.contact_group_02.id
  
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
* ``name`` - (Required) The name for this object.
* ``title`` - (Optional) The title for this object.
* ``phone`` - (Optional) The phone for this object.
* ``email`` - (Optional) The e-mail for this object.
* ``address`` - (Optional) The address for this object.
* ``comments`` - (Optional) Comments for this object.
* ``contact_group_id`` - (Optional) ID of the group where this object belongs to.

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
* ``content_type`` - The content type of this object.

## Import

Tenants can be imported by `id` e.g.

```
$ terraform import netbox_tenancy_contact.contact_test id
```
