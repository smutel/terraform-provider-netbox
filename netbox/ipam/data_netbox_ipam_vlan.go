package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func DataNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about vlan (ipam module) from netbox.",
		ReadContext: dataNetboxIpamVlanRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vlan (ipam module).",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the vlan (ipam module).",
			},
			"vlan_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vlan group where this vlan is attached to.",
			},
		},
	}
}

func dataNetboxIpamVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	vlanID := []int32{int32(d.Get("vlan_id").(int))}
	groupID := int32(d.Get("vlan_group_id").(int))
	groupIDArray := []*int32{&groupID}

	request := client.IpamAPI.IpamVlansList(ctx).Vid(vlanID)
	request = request.GroupId(groupIDArray)

	resource, response, err := request.Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil, errors.New("Your query returned no results. "+
			"Please change your search criteria and try again."))

	} else if resource.GetCount() > 1 {
		return util.GenerateErrorMessage(nil, errors.New("Your query returned more than one result. "+
			"Please try a more specific search criteria."))
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type", util.ConvertURLContentType(r.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
