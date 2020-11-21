# netbox\_json\_ipam\_rirs\_list Data Source

Get json output from the ipam_rirs_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_ipam_rirs_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_ipam_rirs_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

