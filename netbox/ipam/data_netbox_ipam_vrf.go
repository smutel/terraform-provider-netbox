package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func DataNetboxIpamVrf() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about vrf (ipam module) from netbox.",
		ReadContext: dataNetboxIpamVrfRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vrf (ipam module).",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the vrf (ipam module).",
			},
		},
	}
}

func dataNetboxIpamVrfRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	id := int64(d.Get("vrf_id").(int))
	idStr := strconv.FormatInt(id, 10)

	p := ipam.NewIpamVrfsListParams().WithID(&idStr)

	list, err := client.Ipam.IpamVrfsList(p, nil)
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
	if err = d.Set("content_type", util.ConvertURIContentType(r.URL)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
