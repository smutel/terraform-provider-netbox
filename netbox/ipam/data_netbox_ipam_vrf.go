package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netbox "github.com/netbox-community/go-netbox/v4"
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
	client := m.(*netbox.APIClient)

	id := []int32{d.Get("vrf_id").(int32)}

	resource, response, err := client.IpamAPI.IpamVrfsList(ctx).Id(id).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil, errors.New("Your query returned no results. "+
			"Please change your search criteria and try again."))

	} else if resource.GetCount() > 1 {
		return util.GenerateErrorMessage(nil, errors.New("Your query returned more than one result. "+
			"Please try a more specific search criteria."))
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type", util.ConvertURLContentType(r.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
