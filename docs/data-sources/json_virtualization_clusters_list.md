# netbox\_json\_virtualization\_clusters\_list Data Source

Get json output from the virtualization_clusters_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_virtualization_clusters_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_virtualization_clusters_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

