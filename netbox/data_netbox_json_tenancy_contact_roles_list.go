package netbox

import (
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  netboxclient "github.com/smutel/go-netbox/netbox/client"
  "github.com/smutel/go-netbox/netbox/client/tenancy"
)

func dataNetboxJSONTenancyContactRolesList() *schema.Resource {
  return &schema.Resource{
    Read: dataNetboxJSONTenancyContactRolesListRead,

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

func dataNetboxJSONTenancyContactRolesListRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*netboxclient.NetBoxAPI)

  params := tenancy.NewTenancyContactRolesListParams()
  limit := int64(d.Get("limit").(int))
  params.Limit = &limit

  list, err := client.Tenancy.TenancyContactRolesList(params, nil)
  if err != nil {
    return err
  }

  j, _ := json.Marshal(list.Payload.Results)

  d.Set("json", string(j))
  d.SetId("NetboxJSONTenancyContactRolesList")

  return nil
}
