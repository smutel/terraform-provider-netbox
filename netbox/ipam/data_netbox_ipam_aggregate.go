package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func DataNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about aggregate (ipam module) from Netbox.",
		ReadContext: dataNetboxIpamAggregateRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this aggregate (ipam module).",
			},
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
				Description:  "The prefix (with mask) used for this aggregate (ipam module).",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The RIR id linked to this aggregate (ipam module).",
			},
		},
	}
}

func dataNetboxIpamAggregateRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	prefix := d.Get("prefix").(string)
	rirID := d.Get("rir_id").(string)

	p := ipam.NewIpamAggregatesListParams().WithPrefix(&prefix).WithRirID(&rirID)

	list, err := client.Ipam.IpamAggregatesList(p, nil)
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
