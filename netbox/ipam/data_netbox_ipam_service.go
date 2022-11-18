package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func DataNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about a service (ipam module) from netbox.",
		ReadContext: dataNetboxIpamServiceRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this service (ipam module).",
			},
			"device_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"virtualmachine_id"},
				Description:   "ID of the device linked to this service (ipam module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name of this service (ipam module).",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The port of this service (ipam module).",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				Description:  "The protocol of this service (ipam module) (tcp or udp).",
			},
			"virtualmachine_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"device_id"},
				Description:   "ID of the VM linked to this service (ipam module).",
			},
		},
	}
}

func dataNetboxIpamServiceRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	deviceID := int64(d.Get("device_id").(int))
	deviceIDStr := strconv.FormatInt(deviceID, 10)
	name := d.Get("name").(string)
	port := float64(d.Get("port").(int))
	protocol := d.Get("protocol").(string)
	vmID := int64(d.Get("virtualmachine_id").(int))
	vmIDStr := strconv.FormatInt(vmID, 10)

	p := ipam.NewIpamServicesListParams().WithName(&name)
	p.SetPort(&port)
	p.SetProtocol(&protocol)
	if deviceID != 0 {
		p.SetDeviceID(&deviceIDStr)
	} else if vmID != 0 {
		p.SetVirtualMachineID(&vmIDStr)
	}

	list, err := client.Ipam.IpamServicesList(p, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if *list.Payload.Count < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	} else if *list.Payload.Count > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	r := list.Payload.Results[0]
	d.SetId(strconv.FormatInt(r.ID, 10))
	if err = d.Set("content_type", util.ConvertURIContentType(r.URL)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
