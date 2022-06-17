package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/extras"
)

func dataNetboxJSONExtrasTagsList() *schema.Resource {
  return &schema.Resource{
    Description: "Get json output from the extras_tags_list Netbox endpoint.",
    Read: dataNetboxJSONExtrasTagsListRead,

    Schema: map[string]*schema.Schema{
      "limit": {
        Type:     schema.TypeInt,
        Optional: true,
        Default:  0,
        Description: "The max number of returned results. If 0 is specified, all records will be returned.",
      },
      "json": {
        Type:     schema.TypeString,
        Computed: true,
        Description: "JSON output of the list of objects for this Netbox endpoint.",
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
