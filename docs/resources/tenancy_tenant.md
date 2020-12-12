# netbox\_tenancy\_tenant Resource

Manage a tenant within Netbox.

## Example Usage

```hcl
resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "TestTenant"
  slug            = "TestTenant"
  description     = "Tenant created by terraform"
  comments        = "Some test comments"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id
  
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
* ``comments`` - (Optional) Comments for this object.
* ``custom_fields`` - (Optional) Custom Field Keys and Values for this object
  * For boolean, use the string value "true" or "false"
  * For data, use the string format "YYYY-MM-DD"
  * For integer, use the value between double quote "10"
  * For selection, use the level id
  * For text, use the string value
  * For URL, use the URL as string
* ``description`` - (Optional) The description for this object.
* ``tenant_group_id`` - (Optional) ID of the group where this object is located.
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.

The ``tag`` block supports:
* ``name`` - (Required) Name of the existing tag to associate with this resource.
* ``slug`` - (Required) Slug of the existing tag to associate with this resource.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
