# netbox\_json\_extras\_tags\_list Data Source

Get json output from the extras_tags_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_extras_tags_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_extras_tags_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

