package netbox

import (
        "encoding/json"

        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        netboxclient "github.com/netbox-community/go-netbox/netbox/client"
)

func dataNetboxJSONDcimPlatformsList() *schema.Resource {
        return &schema.Resource{
                Read: dataNetboxJSONDcimPlatformsListRead,

                Schema: map[string]*schema.Schema{
                        "json": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                },
        }
}

func dataNetboxJSONDcimPlatformsListRead(d *schema.ResourceData, m interface{}) error {
        client := m.(*netboxclient.NetBoxAPI)

        list, err := client.Dcim.DcimPlatformsList(nil, nil)
        if err != nil {
                return err
        }

        j, _ := json.Marshal(list.Payload.Results)

        d.Set("json", string(j))
        d.SetId("NetboxJSONDcimPlatformsList")

        return nil
}
