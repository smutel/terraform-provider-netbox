package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

func dataNetboxJSONDcimPowerConnectionsList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONDcimPowerConnectionsListRead,

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

func dataNetboxJSONDcimPowerConnectionsListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := dcim.NewDcimPowerConnectionsListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Dcim.DcimPowerConnectionsList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONDcimPowerConnectionsList")

	return nil
}
