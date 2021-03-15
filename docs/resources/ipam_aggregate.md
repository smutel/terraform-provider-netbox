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
* ``custom_fields`` - (Optional) Custom Field Keys and Values for this object
  * For boolean, use the string value "true" or "false"
  * For data, use the string format "YYYY-MM-DD"
  * For integer, use the value between double quote "10"
  * For selection, use the level id
  * For text, use the string value
  * For URL, use the URL as string
* ``date_added`` - (Optional) Date when this aggregate was added. Format *YYYY-MM-DD*.
* ``description`` - (Optional) The description of this object.
* ``prefix`` - (Required) The prefix (with mask) used for this object.
* ``rir_id`` - (Required) The RIR id linked to this object.

The ``tag`` block supports:
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
