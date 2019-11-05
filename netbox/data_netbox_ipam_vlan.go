package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	pkgerrors "github.com/pkg/errors"
)

func dataNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxIpamVlanRead,

		Schema: map[string]*schema.Schema{
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
	client := m.(*netboxclient.NetBox)

	vlanVid := int64(d.Get("vlan_id").(int))
	vlanVidStr := strconv.FormatInt(vlanVid, 10)
	vlanGroupID := int64(d.Get("vlan_group_id").(int))
	vlanGroupIDStr := strconv.FormatInt(vlanGroupID, 10)

	p := ipam.NewIpamVlansListParams().WithVid(&vlanVidStr)
	if vlanGroupID != 0 {
		p.SetGroupID(&vlanGroupIDStr)
	}

	vlanGroupsList, err := client.Ipam.IpamVlansList(p, nil)
	if err != nil {
		return err
	}

	if *vlanGroupsList.Payload.Count == 1 {
		d.SetId(strconv.FormatInt(vlanGroupsList.Payload.Results[0].ID, 10))
	} else {
		return pkgerrors.New("Data results for netbox_ipam_vlan returns 0 or " +
			"more than one result.")
	}

	return nil
}
