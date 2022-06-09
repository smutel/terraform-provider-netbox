package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/ipam"
)

func dataNetboxJSONIpamVlansList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONIpamVlansListRead,

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

func dataNetboxJSONIpamVlansListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := ipam.NewIpamVlansListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Ipam.IpamVlansList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONIpamVlansList")

  return nil
}
