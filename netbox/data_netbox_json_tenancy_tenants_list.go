package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/tenancy"
)

func dataNetboxJSONTenancyTenantsList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONTenancyTenantsListRead,

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

func dataNetboxJSONTenancyTenantsListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := tenancy.NewTenancyTenantsListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Tenancy.TenancyTenantsList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONTenancyTenantsList")

	return nil
}
