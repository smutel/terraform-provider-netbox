package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about a service from netbox.",
		ReadContext: dataNetboxIpamServiceRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this service.",
			},
			"device_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"virtualmachine_id"},
				Description:   "ID of the device linked to this service.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name of this service.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, util.Const65535),
				Description:  "The port of this service.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"},
					false),
				Description: "The protocol of this service (tcp or udp).",
			},
			"virtualmachine_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"device_id"},
				Description:   "ID of the VM linked to this service.",
			},
		},
	}
}

func dataNetboxIpamServiceRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	deviceID32, err := safecast.ToInt32(d.Get("device_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	deviceIDArray := []*int32{&deviceID32}

	name := []string{d.Get("name").(string)}
	port := float32(d.Get("port").(int))
	protocol := d.Get("protocol").(string)

	vmID32, err := safecast.ToInt32(d.Get("virtualmachine_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	vmIDArray := []*int32{&vmID32}

	request := client.IpamAPI.IpamServicesList(ctx).Name(name)
	request = request.Port(port)
	p, err :=
		netbox.NewIpamServiceTemplatesListProtocolParameterFromValue(
			protocol)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	request = request.Protocol(*p)

	if deviceID32 != 0 {
		request = request.DeviceId(deviceIDArray)
	} else if vmID32 != 0 {
		request = request.VirtualMachineId(vmIDArray)
	}

	resource, response, err := request.Execute()

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
