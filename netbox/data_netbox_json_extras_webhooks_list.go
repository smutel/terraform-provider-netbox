package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/extras"
)

func dataNetboxJSONExtrasWebhooksList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONExtrasWebhooksListRead,

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

func dataNetboxJSONExtrasWebhooksListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := extras.NewExtrasWebhooksListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Extras.ExtrasWebhooksList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONExtrasWebhooksList")

  return nil
}
