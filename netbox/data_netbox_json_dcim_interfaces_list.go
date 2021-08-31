package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

func dataNetboxJSONDcimInterfacesList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONDcimInterfacesListRead,

		Schema: map[string]*schema.Schema{
			"devicename": {
				Type:     schema.TypeString,
				Optional: true,
			},
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

func dataNetboxJSONDcimInterfacesListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	deviceName := d.Get("devicename").(string)
	limit := int64(d.Get("limit").(int))
	params := dcim.NewDcimInterfacesListParams()
	if deviceName != "" {
		params.Device = &deviceName
	}
	params.Limit = &limit

	list, err := client.Dcim.DcimInterfacesList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONDcimInterfacesList")

	return nil
}
