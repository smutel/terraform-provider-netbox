package netbox

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/tenancy"
	pkgerrors "github.com/pkg/errors"
)

func dataNetboxTenancyTenant() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxTenancyTenantRead,

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

func dataNetboxTenancyTenantRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBox)

	slug := d.Get("slug").(string)

	p := tenancy.NewTenancyTenantsListParams().WithSlug(&slug)

	list, err := client.Tenancy.TenancyTenantsList(p, nil)
	if err != nil {
		return err
	}

	if *list.Payload.Count == 1 {
		d.SetId(strconv.FormatInt(list.Payload.Results[0].ID, 10))
	} else {
		return pkgerrors.New("Data results for netbox_tenancy_tenant returns 0 " +
			"or more than one result.")
	}

	return nil
}
