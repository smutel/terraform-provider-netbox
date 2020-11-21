# netbox\_json\_dcim\_power\_ports\_list Data Source

Get json output from the dcim_power_ports_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_dcim_power_ports_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_dcim_power_ports_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

