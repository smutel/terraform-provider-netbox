package netbox

import (
        "encoding/json"

        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        netboxclient "github.com/netbox-community/go-netbox/netbox/client"
)

func dataNetboxJSONDcimDevicesList() *schema.Resource {
        return &schema.Resource{
                Read: dataNetboxJSONDcimDevicesListRead,

                Schema: map[string]*schema.Schema{
                        "json": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                },
        }
}

func dataNetboxJSONDcimDevicesListRead(d *schema.ResourceData, m interface{}) error {
        client := m.(*netboxclient.NetBoxAPI)

        list, err := client.Dcim.DcimDevicesList(nil, nil)
        if err != nil {
                return err
        }

        j, _ := json.Marshal(list.Payload.Results)

        d.Set("json", string(j))
        d.SetId("NetboxJSONDcimDevicesList")

        return nil
}
