package netbox

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/tenancy"
)

func dataNetboxTenancyContactGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataNetboxTenancyContactGroupRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
	}
}

func dataNetboxTenancyContactGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	slug := d.Get("slug").(string)

	p := tenancy.NewTenancyContactGroupsListParams().WithSlug(&slug)

	list, err := client.Tenancy.TenancyContactGroupsList(p, nil)
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
	d.Set("content_type", convertURIContentType(r.URL))

	return nil
}
