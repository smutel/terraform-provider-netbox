# netbox\_json\_extras\_image\_attachments\_list Data Source

Get json output from the extras_image_attachments_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_extras_image_attachments_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_extras_image_attachments_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

