# netbox\_json\_dcim\_device\_bay\_templates\_list Data Source

Get json output from the dcim_device_bay_templates_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_dcim_device_bay_templates_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_dcim_device_bay_templates_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

