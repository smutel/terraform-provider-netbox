#!/bin/bash

FILE=netbox_list_endpoints.txt

while read -r line; do
        SPLIT=( $line )

        ENDPOINT=${SPLIT[0]}
        SECTION=${SPLIT[1]}
        ITEM=${SPLIT[2]}

cat << EOF > ../netbox/data_netbox_json_${ENDPOINT}_list.go
package netbox

import (
        "encoding/json"

        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        netboxclient "github.com/netbox-community/go-netbox/netbox/client"
)

func dataNetboxJson${SECTION}${ITEM}List() *schema.Resource {
        return &schema.Resource{
                Read: dataNetboxJson${SECTION}${ITEM}ListRead,

                Schema: map[string]*schema.Schema{
                        "json": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                },
        }
}

func dataNetboxJson${SECTION}${ITEM}ListRead(d *schema.ResourceData, m interface{}) error {
        client := m.(*netboxclient.NetBoxAPI)

        list, err := client.${SECTION}.${SECTION}${ITEM}List(nil, nil)
        if err != nil {
                return err
        }

        j, _ := json.Marshal(list.Payload.Results)

        d.Set("json", string(j))
        d.SetId("NetboxJson${SECTION}${ITEM}List")

        return nil
}
EOF


cat << EOF >> provider_update.txt
"netbox_json_${ENDPOINT}_list": dataNetboxJson${SECTION}${ITEM}List(),
EOF


cat << EOF > ../docs/data-sources/json_${ENDPOINT}_list.md
# netbox\_json\_`echo ${ENDPOINT} | sed 's/_/\\\_/g'`\_list Data Source

Get json output from the ${ENDPOINT}_list Netbox endpoint

## Example Usage

\`\`\`hcl
data "netbox_json_${ENDPOINT}_list" "test" {}
output "example" {
  value = jsondecode(data.netbox_json_${ENDPOINT}_list.test.json)
}
\`\`\`

## Argument Reference

This function takes no arguments.

## Attributes Reference

This function has no additional attributes.

EOF

done < "$FILE"
