package netbox

import (
        "encoding/json"

        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        netboxclient "github.com/netbox-community/go-netbox/netbox/client"
)

func dataNetboxJSONIpamVlansList() *schema.Resource {
        return &schema.Resource{
                Read: dataNetboxJSONIpamVlansListRead,

                Schema: map[string]*schema.Schema{
                        "json": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                },
        }
}

func dataNetboxJSONIpamVlansListRead(d *schema.ResourceData, m interface{}) error {
        client := m.(*netboxclient.NetBoxAPI)

        list, err := client.Ipam.IpamVlansList(nil, nil)
        if err != nil {
                return err
        }

        j, _ := json.Marshal(list.Payload.Results)

        d.Set("json", string(j))
        d.SetId("NetboxJSONIpamVlansList")

        return nil
}
