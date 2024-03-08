package dcim

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func DataNetboxDcimLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about location (dcim module) from netbox.",
		ReadContext: dataNetboxDcimLocationRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this location (dcim module).",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug of the location (dcim module).",
			},
		},
	}
}

func dataNetboxDcimLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	slug := []string{d.Get("slug").(string)}

	resource, response, err := client.DcimAPI.DcimLocationsList(ctx).Slug(slug).Execute()

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
