# netbox\_json\_extras\_job\_results\_list Data Source

Get json output from the extras_job_results_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_extras_job_results_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_extras_job_results_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

