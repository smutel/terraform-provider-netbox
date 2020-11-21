# netbox\_json\_virtualization\_interfaces\_list Data Source

Get json output from the virtualization_interfaces_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_virtualization_interfaces_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_virtualization_interfaces_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

