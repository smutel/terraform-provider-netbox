# netbox\_json\_ipam\_services\_list Data Source

Get json output from the ipam_services_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_ipam_services_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_ipam_services_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

