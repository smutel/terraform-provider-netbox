package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func DataNetboxIpamIPAddresses() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about IP addresses (ipam module) from netbox.",
		ReadContext: dataNetboxIpamIPAddressesRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this ipam IP addresses (ipam module).",
			},
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The address (with mask) of the ipam IP addresses (ipam module).",
			},
		},
	}
}

func dataNetboxIpamIPAddressesRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	address := d.Get("address").(string)

	p := ipam.NewIpamIPAddressesListParams().WithAddress(&address)

	list, err := client.Ipam.IpamIPAddressesList(p, nil)
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
