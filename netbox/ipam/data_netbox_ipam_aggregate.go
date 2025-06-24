package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about aggregate from Netbox.",
		ReadContext: dataNetboxIpamAggregateRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this aggregate.",
			},
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, util.Const256),
				Description:  "The prefix (with mask) used for this aggregate.",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The RIR id linked to this aggregate.",
			},
		},
	}
}

func dataNetboxIpamAggregateRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	prefix := d.Get("prefix").(string)

	rirID32, err := safecast.ToInt32(d.Get("rir_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	rirIDArray := []int32{rirID32}

	resource, response, err :=
		client.IpamAPI.IpamAggregatesList(
			ctx).Prefix(prefix).RirId(rirIDArray).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned no results. "+
				"Please change your search criteria and try again."))

	} else if resource.GetCount() > 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned more than one result. "+
				"Please try a more specific search criteria."))
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type",
		util.ConvertURLContentType(r.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
