# netbox\_json\_users\_permissions\_list Data Source

Get json output from the users_permissions_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_users_permissions_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_users_permissions_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

