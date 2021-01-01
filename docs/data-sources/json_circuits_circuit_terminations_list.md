# netbox\_json\_circuits\_circuit\_terminations\_list Data Source

Get json output from the circuits_circuit_terminations_list Netbox endpoint

## Example Usage

```hcl
data "netbox_json_circuits_circuit_terminations_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_circuits_circuit_terminations_list.test.json)
}
```

## Argument Reference

This function takes no arguments.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* ``json`` - JSON output of the list of objects for this Netbox endpoint.

