package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/dcim"
)

func dataNetboxJSONDcimSitesList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONDcimSitesListRead,

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

func dataNetboxJSONDcimSitesListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := dcim.NewDcimSitesListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Dcim.DcimSitesList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONDcimSitesList")

  return nil
}
