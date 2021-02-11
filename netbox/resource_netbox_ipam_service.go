package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamServiceCreate,
		Read:   resourceNetboxIpamServiceRead,
		Update: resourceNetboxIpamServiceUpdate,
		Delete: resourceNetboxIpamServiceDelete,
		Exists: resourceNetboxIpamServiceExists,

		Schema: map[string]*schema.Schema{
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// terraform default behavior sees a difference between null and an empty string
				// therefore we override the default, because null in terraform results in empty string in netbox
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// function is called for each member of map
					// including additional call on the amount of entries
					// we ignore the count, because the actual state always returns the amount of existing custom_fields and all are optional in terraform
					if k == CustomFieldsRegex {
						return true
					}
					return old == new
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      " ",
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"device_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"device_id", "virtualmachine_id"},
			},
			"ip_addresses_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"virtualmachine_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNetboxIpamServiceCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_fields").(map[string]interface{})
	customFields := convertCustomFieldsFromTerraformToAPICreate(resourceCustomFields)
	description := d.Get("description").(string)
	deviceID := int64(d.Get("device_id").(int))
	IPaddressesID := d.Get("ip_addresses_id").([]interface{})
	IPaddressesID64 := []int64{}
	name := d.Get("name").(string)
	port := int64(d.Get("port").(int))
	protocol := d.Get("protocol").(string)
	tags := d.Get("tag").(*schema.Set).List()
	virtualmachineID := int64(d.Get("virtualmachine_id").(int))

	for _, id := range IPaddressesID {
		IPaddressesID64 = append(IPaddressesID64, int64(id.(int)))
	}

	newResource := &models.WritableService{
		CustomFields: &customFields,
		Description:  description,
		Ipaddresses:  IPaddressesID64,
		Name:         &name,
		Port:         &port,
		Protocol:     &protocol,
		Tags:         convertTagsToNestedTags(tags),
	}

	if deviceID != 0 {
		newResource.Device = &deviceID
	}

	if virtualmachineID != 0 {
		newResource.VirtualMachine = &virtualmachineID
	}

	resource := ipam.NewIpamServicesCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamServicesCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamServiceRead(d, m)
}

func resourceNetboxIpamServiceRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamServicesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamServicesList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			customFields := convertCustomFieldsFromAPIToTerraform(resource.CustomFields)

			if err = d.Set("custom_fields", customFields); err != nil {
				return err
			}

			var description string

			if resource.Description == "" {
				description = " "
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return err
			}

			if resource.Device == nil {
				if err = d.Set("device_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("device_id", resource.Device.ID); err != nil {
					return err
				}
			}

			IPaddressesObject := resource.Ipaddresses
			IPaddressesInt := []int64{}
			for _, ip := range IPaddressesObject {
				IPaddressesInt = append(IPaddressesInt, ip.ID)
			}
			if err = d.Set("ip_addresses_id", IPaddressesInt); err != nil {
				return err
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if err = d.Set("port", resource.Port); err != nil {
				return err
			}

			if err = d.Set("protocol", resource.Protocol.Value); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			if resource.VirtualMachine == nil {
				if err = d.Set("virtualmachine_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("virtualmachine_id", resource.VirtualMachine.ID); err != nil {
					return err
				}
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamServiceUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableService{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	port := int64(d.Get("port").(int))
	params.Port = &port

	protocol := d.Get("protocol").(string)
	params.Protocol = &protocol

	IPaddressesID := d.Get("ip_addresses_id").([]interface{})
	IPaddressesID64 := []int64{}
	for _, id := range IPaddressesID {
		IPaddressesID64 = append(IPaddressesID64, int64(id.(int)))
	}

	params.Ipaddresses = IPaddressesID64

	deviceID := int64(d.Get("device_id").(int))
	if deviceID != 0 {
		params.Device = &deviceID
	}

	vmID := int64(d.Get("virtualmachine_id").(int))
	if vmID != 0 {
		params.VirtualMachine = &vmID
	}

	if d.HasChange("custom_fields") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_fields")
		customFields := convertCustomFieldsFromTerraformToAPIUpdate(stateCustomFields, resourceCustomFields)
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	resource := ipam.NewIpamServicesPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamServicesPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamServiceRead(d, m)
}

func resourceNetboxIpamServiceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamServiceExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource := ipam.NewIpamServicesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamServicesDelete(resource, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxIpamServiceExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	vlanID := d.Id()
	params := ipam.NewIpamServicesListParams().WithID(&vlanID)
	vlans, err := client.Ipam.IpamServicesList(params, nil)

	if err != nil {
		return resourceExist, err
	}

	for _, vlan := range vlans.Payload.Results {
		if strconv.FormatInt(vlan.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
