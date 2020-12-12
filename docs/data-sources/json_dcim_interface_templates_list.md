# netbox\_json\_dcim\_interface\_templates\_list Data Source

Get json output from the dcim_interface_templates_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_dcim_interface_templates_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_dcim_interface_templates_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.
