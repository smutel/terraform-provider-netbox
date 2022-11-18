package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func DataNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about vlan (ipam module) from netbox.",
		ReadContext: dataNetboxIpamVlanRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vlan (ipam module).",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the vlan (ipam module).",
			},
			"vlan_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vlan group where this vlan is attached to.",
			},
		},
	}
}

func dataNetboxIpamVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	if *list.Payload.Count < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	} else if *list.Payload.Count > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	r := list.Payload.Results[0]
	d.SetId(strconv.FormatInt(r.ID, 10))
	d.Set("content_type", util.ConvertURIContentType(r.URL))

	return nil
}
