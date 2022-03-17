package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/extras"
)

func dataNetboxJSONExtrasTagsList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONExtrasTagsListRead,

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

func dataNetboxJSONExtrasTagsListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := extras.NewExtrasTagsListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Extras.ExtrasTagsList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONExtrasTagsList")

  return nil
}
