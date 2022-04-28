package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamServiceCreate,
		Read:   resourceNetboxIpamServiceRead,
		Update: resourceNetboxIpamServiceUpdate,
		Delete: resourceNetboxIpamServiceDelete,
		Exists: resourceNetboxIpamServiceExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
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
			"ports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	deviceID := int64(d.Get("device_id").(int))
	IPaddressesID := d.Get("ip_addresses_id").([]interface{})
	IPaddressesID64 := []int64{}
	name := d.Get("name").(string)
	ports := d.Get("ports").([]interface{})
	ports64 := []int64{}
	protocol := d.Get("protocol").(string)
	tags := d.Get("tag").(*schema.Set).List()
	virtualmachineID := int64(d.Get("virtualmachine_id").(int))

	for _, id := range IPaddressesID {
		IPaddressesID64 = append(IPaddressesID64, int64(id.(int)))
	}

	for _, p := range ports {
		ports64 = append(ports64, int64(p.(int)))
	}

	newResource := &models.WritableService{
		CustomFields: &customFields,
		Description:  description,
		Ipaddresses:  IPaddressesID64,
		Name:         &name,
		Ports:        ports64,
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
			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return err
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
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

			ports := resource.Ports
			if err = d.Set("ports", ports); err != nil {
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

	ports := d.Get("ports").([]interface{})
	ports64 := []int64{}
	for _, port := range ports {
		ports64 = append(ports64, int64(port.(int)))
	}

	params.Ports = ports64

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

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
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
