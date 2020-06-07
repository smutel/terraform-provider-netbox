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

func dataNetboxIpamIPAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxIpamIPAddressesRead,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}/"+
						"[0-9]{1,2}$"), "Must be like 192.168.56.1/24"),
			},
		},
	}
}

func dataNetboxIpamIPAddressesRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	address := d.Get("address").(string)

	p := ipam.NewIpamIPAddressesListParams().WithAddress(&address)

	list, err := client.Ipam.IpamIPAddressesList(p, nil)
	if err != nil {
		return err
	}

	if *list.Payload.Count == 1 {
		d.SetId(strconv.FormatInt(list.Payload.Results[0].ID, 10))
	} else {
		return pkgerrors.New("Data results for netbox_ipam_ip_addresses returns 0 or " +
			"more than one result.")
	}

	return nil
}
