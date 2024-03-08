package ipam

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
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
	client := m.(*netbox.APIClient)

	prefix := d.Get("prefix").(string)
	rirID := []int32{int32(d.Get("rir_id").(int))}

	resource, response, err := client.IpamAPI.IpamAggregatesList(ctx).Prefix(prefix).RirId(rirID).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	} else if resource.GetCount() > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type", util.ConvertURIContentType(strfmt.URI(r.GetUrl()))); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
