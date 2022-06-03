package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/virtualization"
)

func dataNetboxVirtualizationCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxVirtualizationClusterRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
	}
}

func dataNetboxVirtualizationClusterRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	name := d.Get("name").(string)

	resource := virtualization.NewVirtualizationClustersListParams().WithName(
		&name)

	list, err := client.Virtualization.VirtualizationClustersList(resource, nil)
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

	r := list.Payload.Results[0]
	d.SetId(strconv.FormatInt(r.ID, 10))
	d.Set("content_type", convertURIContentType(r.URL))

	return nil
}
