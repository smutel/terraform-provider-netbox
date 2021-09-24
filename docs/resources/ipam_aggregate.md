# netbox\_ipam\_aggregate Resource

Manage an aggregate within Netbox.

## Example Usage

```hcl
resource "netbox_ipam_aggregate" "aggregate_test" {
  prefix = "192.168.56.0/24"
  rir_id = 1
  date_created = "2020-12-21"
  
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
* ``date_added`` - (Optional) Date when this aggregate was added. Format *YYYY-MM-DD*.
* ``description`` - (Optional) The description of this object.
* ``prefix`` - (Required) The prefix (with mask) used for this object.
* ``rir_id`` - (Required) The RIR id linked to this object.

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

Aggregates can be imported by `id` e.g.

```
$ terraform import netbox_ipam_aggregate.aggregate_test id
```
