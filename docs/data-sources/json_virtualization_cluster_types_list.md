# netbox\_json\_virtualization\_cluster\_types\_list Data Source

Get json output from the virtualization_cluster_types_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_virtualization_cluster_types_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_virtualization_cluster_types_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

