package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/circuits"
)

func dataNetboxJSONCircuitsProviderNetworksList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONCircuitsProviderNetworksListRead,

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

func dataNetboxJSONCircuitsProviderNetworksListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := circuits.NewCircuitsProviderNetworksListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Circuits.CircuitsProviderNetworksList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONCircuitsProviderNetworksList")

  return nil
}
