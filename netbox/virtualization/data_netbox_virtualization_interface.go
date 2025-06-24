package virtualization

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxVirtualizationInterface() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about interface from Netbox.",
		ReadContext: dataNetboxVirtualizationInterfaceRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this interface.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const64),
				Description:  "The name of this interface.",
			},
			"virtualmachine_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "ID of the VM where this interface " +
					"is attached to.",
			},
		},
	}
}

func dataNetboxVirtualizationInterfaceRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	name := []string{d.Get("name").(string)}
	vmID32, err := safecast.ToInt32(d.Get("virtualmachine_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	vmID := []int32{vmID32}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationInterfacesList(
			ctx).Name(name).VirtualMachineId(vmID).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned no results. "+
				"Please change your search criteria and try again."))

	} else if resource.GetCount() > 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned more than one result. "+
				"Please try a more specific search criteria."))
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type",
		util.ConvertURLContentType(r.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
