package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
)

func dataNetboxJSONVirtualizationClustersList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONVirtualizationClustersListRead,

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

func dataNetboxJSONVirtualizationClustersListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := virtualization.NewVirtualizationClustersListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Virtualization.VirtualizationClustersList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONVirtualizationClustersList")

	return nil
}
