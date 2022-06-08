# netbox\_virtualization\_cluster Data Source

Get info about vitualization cluster from netbox.

## Example Usage

```hcl
data "netbox_virtualization_cluster" "cluster_test" {
  name = "TestCluster"
}
```

## Argument Reference

The following arguments are supported:
* ``name`` - (Required) The name of the cluster.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
* ``content_type`` - The content type of this object.
