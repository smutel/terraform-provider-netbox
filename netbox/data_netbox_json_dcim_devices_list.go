package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

func dataNetboxJSONDcimDevicesList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONDcimDevicesListRead,

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

func dataNetboxJSONDcimDevicesListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := dcim.NewDcimDevicesListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Dcim.DcimDevicesList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONDcimDevicesList")

	return nil
}
