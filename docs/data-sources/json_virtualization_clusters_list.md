# netbox\_json\_virtualization\_clusters\_list Data Source

Get json output from the virtualization_clusters_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_virtualization_clusters_list" "test" {
  filter {
    name = ""
    value = ""
  }
  limit = 0
}
output "example" {
  value = jsondecode(data.netbox_json_virtualization_clusters_list.test.json)
}
```

## Argument Reference

The following arguments are supported:

* ``filter`` (Optional). This block set should include "name" (String) and "value" (String).
  Supported options:
  - ``name`` The name of the cluster.
  - ``name_ic`` The string containing some part of the cluster name.
  - ``tag`` The tag of a cluster.
* ``limit`` (Optional). The max number of returned results. If 0 is specified, all clusters will be returned.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

