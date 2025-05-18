package virtualization

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxVirtualizationCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about cluster from Netbox.",
		ReadContext: dataNetboxVirtualizationClusterRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:				schema.TypeString,
				Computed:		true,
				Description: "The content type of this cluster.",
			},
			"name": {
				Type:				 schema.TypeString,
				Required:		 true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:	"The name of this cluster.",
			},
			"site_id": {
				Type:				schema.TypeInt,
				Computed:		true,
				Description: "The ID of site where this cluster is located.",
			},
		},
	}
}

func dataNetboxVirtualizationClusterRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	name := []string{d.Get("name").(string)}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationClustersList(ctx).Name(
			name).Execute()

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

	if err = d.Set("site_id", r.GetSite().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
