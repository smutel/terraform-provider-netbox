package netbox

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	pkgerrors "github.com/pkg/errors"
)

func dataNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxIpamVlanGroupRead,

		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataNetboxIpamVlanGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBox)

	slug := d.Get("slug").(string)
	siteID := d.Get("site_id").(int)
	siteIDStr := strconv.FormatInt(int64(siteID), 10)

	p := ipam.NewIpamVlanGroupsListParams().WithSlug(&slug)
	if siteID != 0 {
		p.SetSiteID(&siteIDStr)
	}

	list, err := client.Ipam.IpamVlanGroupsList(p, nil)
	if err != nil {
		return err
	}

	if *list.Payload.Count == 1 {
		d.SetId(strconv.FormatInt(list.Payload.Results[0].ID, 10))
	} else {
		return pkgerrors.New("Data results for netbox_ipam_vlan_group returns 0 " +
			"or more than one result.")
	}

	return nil
}
