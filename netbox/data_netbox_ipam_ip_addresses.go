package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
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
	client := m.(*netboxclient.NetBoxAPI)

	address := d.Get("address").(string)

	p := ipam.NewIpamIPAddressesListParams().WithAddress(&address)

	list, err := client.Ipam.IpamIPAddressesList(p, nil)
	if err != nil {
		return err
	}

	if *list.Payload.Count < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	} else if *list.Payload.Count > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	d.SetId(strconv.FormatInt(list.Payload.Results[0].ID, 10))

	return nil
}
