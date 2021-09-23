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
        "github.com/netbox-community/go-netbox/netbox/client/${SECTION,}"
)

func dataNetboxJSON${SECTION}${ITEM}List() *schema.Resource {
        return &schema.Resource{
                Read: dataNetboxJSON${SECTION}${ITEM}ListRead,

                Schema: map[string]*schema.Schema{
                        "limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
                        "json": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                },
        }
}

func dataNetboxJSON${SECTION}${ITEM}ListRead(d *schema.ResourceData, m interface{}) error {
        client := m.(*netboxclient.NetBoxAPI)

        params := ${SECTION,}.New${SECTION}${ITEM}ListParams()
        limit := int64(d.Get("limit").(int))
        params.Limit = &limit

        list, err := client.${SECTION}.${SECTION}${ITEM}List(params, nil)
        if err != nil {
                return err
        }

        j, _ := json.Marshal(list.Payload.Results)

        d.Set("json", string(j))
        d.SetId("NetboxJSON${SECTION}${ITEM}List")

        return nil
}
EOF


cat << EOF >> provider_update.txt
"netbox_json_${ENDPOINT}_list": dataNetboxJSON${SECTION}${ITEM}List(),
EOF


cat << EOF > ../docs/data-sources/json_${ENDPOINT}_list.md
# netbox\_json\_`echo ${ENDPOINT} | sed 's/_/\\\_/g'`\_list Data Source

Get json output from the ${ENDPOINT}_list Netbox endpoint

## Example Usage

\`\`\`hcl
data "netbox_json_${ENDPOINT}_list" "test" {
        limit = 0
}
output "example" {
        value = jsondecode(data.netbox_json_${ENDPOINT}_list.test.json)
}
\`\`\`

## Argument Reference

* \`\`limit\`\` (Optional). The max number of returned results. If 0 is specified, all records will be returned.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:
* \`\`json\`\` - JSON output of the list of objects for this Netbox endpoint.

EOF

done < "$FILE"
