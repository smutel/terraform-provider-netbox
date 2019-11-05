package netbox

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	pkgerrors "github.com/pkg/errors"
)

func dataNetboxDcimSite() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxDcimSiteRead,

		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
			},
		},
	}
}

func dataNetboxDcimSiteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBox)

	siteSlug := d.Get("slug").(string)

	p := dcim.NewDcimSitesListParams().WithSlug(&siteSlug)

	sitesList, err := client.Dcim.DcimSitesList(p, nil)
	if err != nil {
		return err
	}

	if *sitesList.Payload.Count == 1 {
		d.SetId(strconv.FormatInt(sitesList.Payload.Results[0].ID, 10))
	} else {
		return pkgerrors.New("Data results for netbox_dcim_site returns 0 or " +
			"more than one result.")
	}

	return nil
}
