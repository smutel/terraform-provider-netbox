package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
)

func dataNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxIpamServiceRead,

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"virtualmachine_id"},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
			},
			"virtualmachine_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"device_id"},
			},
		},
	}
}

func dataNetboxIpamServiceRead(d *schema.ResourceData,
	m interface{}) error {
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
