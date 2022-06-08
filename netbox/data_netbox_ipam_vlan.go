package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
)

func dataNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxIpamVlanRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vlan_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataNetboxIpamVlanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	id := int64(d.Get("vlan_id").(int))
	idStr := strconv.FormatInt(id, 10)
	groupID := int64(d.Get("vlan_group_id").(int))
	groupIDStr := strconv.FormatInt(groupID, 10)

	p := ipam.NewIpamVlansListParams().WithVid(&idStr)
	if groupID != 0 {
		p.SetGroupID(&groupIDStr)
	}

	list, err := client.Ipam.IpamVlansList(p, nil)
	if err != nil {
		return err
	}

	if *list.Payload.Count < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	} else if *list.Payload.Count > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	r := list.Payload.Results[0]
	d.SetId(strconv.FormatInt(r.ID, 10))
	d.Set("content_type", convertURIContentType(r.URL))

	return nil
}
