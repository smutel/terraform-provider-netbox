package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func DataNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about prefix (ipam module) from netbox.",
		ReadContext: dataNetboxIpamPrefixRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this prefix (ipam module).",
			},
			"prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
				Description:  "The prefix (IP address/mask) used for this prefix (ipam module).",
			},
		},
	}
}

func dataNetboxIpamPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	prefix := d.Get("prefix").(string)

	p := ipam.NewIpamPrefixesListParams().WithPrefix(&prefix)

	list, err := client.Ipam.IpamPrefixesList(p, nil)
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
