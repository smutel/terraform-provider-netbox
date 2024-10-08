package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an service (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamServiceCreate,
		ReadContext:   resourceNetboxIpamServiceRead,
		UpdateContext: resourceNetboxIpamServiceUpdate,
		DeleteContext: resourceNetboxIpamServiceDelete,
		Exists:        resourceNetboxIpamServiceExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this service (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this service (ipam module).",
			},
			"device_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"device_id", "virtualmachine_id"},
				Description:  "ID of the device linked to this service (ipam module).",
			},
			"ip_addresses_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array of ID of IP addresses attached to this service (ipam module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name for this service (ipam module).",
			},
			"ports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array of ports of this service (ipam module).",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				Description:  "The protocol of this service (ipam module) (tcp or udp).",
			},
			"tag": &tag.TagSchema,
			"virtualmachine_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the VM linked to this service (ipam module).",
			},
		},
	}
}

func resourceNetboxIpamServiceCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	deviceID := int32(d.Get("device_id").(int))
	IPaddressesID := d.Get("ip_addresses_id").([]int32)
	name := d.Get("name").(string)
	ports := d.Get("ports").([]int32)
	protocol := d.Get("protocol").(string)
	tags := d.Get("tag").(*schema.Set).List()
	virtualmachineID := int32(d.Get("virtualmachine_id").(int))

	newResource := netbox.NewWritableServiceRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetIpaddresses(IPaddressesID)
	newResource.SetName(name)
	newResource.SetPorts(ports)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	p, err := netbox.NewPatchedWritableServiceRequestProtocolFromValue(protocol)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetProtocol(*p)

	if deviceID != 0 {
		newResource.SetDevice(deviceID)
	}

	if virtualmachineID != 0 {
		newResource.SetVirtualMachine(virtualmachineID)
	}

	resourceCreated, response, err := client.IpamAPI.IpamServicesCreate(ctx).WritableServiceRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxIpamServiceRead(ctx, d, m)
}

func resourceNetboxIpamServiceRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.IpamAPI.IpamServicesRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("device_id", resource.GetDevice().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("ip_addresses_id", resource.GetIpaddresses()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("ports", resource.GetPorts()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("protocol", resource.GetProtocol().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vittualmachine_id", resource.GetVirtualMachine().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamServiceUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableServiceRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("device_id") {
		deviceID := int32(d.Get("device_id").(int))
		if deviceID != 0 {
			resource.SetDevice(deviceID)
		} else {
			resource.SetDeviceNil()
		}
	}

	if d.HasChange("ip_addresses_id") {
		IPaddressesID := d.Get("ip_addresses_id").([]int32)
		resource.SetIpaddresses(IPaddressesID)
	}

	if d.HasChange("ports") {
		ports := d.Get("ports").([]int32)
		resource.SetPorts(ports)
	}

	if d.HasChange("protocol") {
		p, err := netbox.NewPatchedWritableServiceRequestProtocolFromValue(d.Get("protocol").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetProtocol(*p)
	}

	if d.HasChange("virtualmachine_id") {
		vmID := int32(d.Get("virtualmachine_id").(int))
		if vmID != 0 {
			resource.SetVirtualMachine(vmID)
		} else {
			resource.SetVirtualMachineNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.IpamAPI.IpamServicesUpdate(ctx, int32(resourceID)).WritableServiceRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamServiceRead(ctx, d, m)
}

func resourceNetboxIpamServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamServiceExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.IpamAPI.IpamServicesDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamServiceExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamServicesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
