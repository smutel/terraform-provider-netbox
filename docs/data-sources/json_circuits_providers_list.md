# netbox\_json\_circuits\_providers\_list Data Source

Get json output from the circuits_providers_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_circuits_providers_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_circuits_providers_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

