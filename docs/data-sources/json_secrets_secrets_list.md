# netbox\_json\_secrets\_secrets\_list Data Source

Get json output from the secrets_secrets_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_secrets_secrets_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_secrets_secrets_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

