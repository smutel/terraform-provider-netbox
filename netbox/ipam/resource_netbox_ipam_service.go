// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamService() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a service within Netbox.",
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
				Description: "The content type of this service.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this VRF was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this service.",
			},
			"device_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"device_id", "virtualmachine_id"},
				Description:  "ID of the device linked to this service.",
			},
			"ip_addresses_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					Description: "One of the IP address " +
						"attached to this service.",
				},
				Description: "Array of ID of IP addresses " +
					"attached to this service.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this service was updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name for this service.",
			},
			"ports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "One of the port for this service.",
				},
				Description: "Array of ports of this service.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"},
					false),
				Description: "The protocol of this service (tcp or udp).",
			},
			"tag": &tag.TagSchema,
			"virtualmachine_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the VM linked to this service.",
			},
		},
	}
}

func resourceNetboxIpamServiceCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	deviceID := d.Get("device_id").(int)

	ipaddresses := d.Get("ip_addresses_id").([]any)
	ipaddressesID, err := util.ExpandToInt32Slice(ipaddresses)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	name := d.Get("name").(string)
	ports := d.Get("ports").([]any)
	portsID, err := util.ExpandToInt32Slice(ports)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	protocol := d.Get("protocol").(string)
	tags := d.Get("tag").(*schema.Set).List()
	virtualmachineID := d.Get("virtualmachine_id").(int)

	newResource := netbox.NewWritableServiceRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetIpaddresses(ipaddressesID)
	newResource.SetName(name)
	newResource.SetPorts(portsID)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	p, err := netbox.NewPatchedWritableServiceRequestProtocolFromValue(
		protocol)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetProtocol(*p)

	if deviceID != 0 {
		b, err := brief.GetBriefDeviceRequestFromID(ctx, client, deviceID)
		if err != nil {
			return err
		}
		newResource.SetDevice(*b)
	}

	if virtualmachineID != 0 {
		b, err := brief.GetBriefVirtualMachineRequestFromID(ctx, client,
			virtualmachineID)
		if err != nil {
			return err
		}
		newResource.SetVirtualMachine(*b)
	}

	_, response, err := client.IpamAPI.IpamServicesCreate(
		ctx).WritableServiceRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamServiceRead(ctx, d, m)
}

func resourceNetboxIpamServiceRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamServicesRetrieve(ctx,
		int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields, resource.GetCustomFields())

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

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
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

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(
		resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_id",
		resource.GetVirtualMachine().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamServiceUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableServiceRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))
	ports := d.Get("ports").([]any)
	portsID, err := util.ExpandToInt32Slice(ports)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	resource.SetPorts(portsID)

	p, err := netbox.NewPatchedWritableServiceRequestProtocolFromValue(
		d.Get("protocol").(string))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	resource.SetProtocol(*p)

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("device_id") {
		deviceID := d.Get("device_id").(int)
		if deviceID != 0 {
			b, err := brief.GetBriefDeviceRequestFromID(ctx, client, deviceID)
			if err != nil {
				return err
			}
			resource.SetDevice(*b)
		} else {
			resource.SetDeviceNil()
		}
	}

	if d.HasChange("ip_addresses_id") {
		ipaddresses := d.Get("ip_addresses_id").([]any)
		ipaddressesID, err := util.ExpandToInt32Slice(ipaddresses)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetIpaddresses(ipaddressesID)
	}

	if d.HasChange("virtualmachine_id") {
		vmID := d.Get("virtualmachine_id").(int)
		if vmID != 0 {
			b, err := brief.GetBriefVirtualMachineRequestFromID(
				ctx, client, vmID)
			if err != nil {
				return err
			}
			resource.SetVirtualMachine(*b)
		} else {
			resource.SetVirtualMachineNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.IpamAPI.IpamServicesUpdate(ctx,
		int32(resourceID)).WritableServiceRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamServiceRead(ctx, d, m)
}

func resourceNetboxIpamServiceDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamServiceExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.IpamAPI.IpamServicesDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamServiceExists(d *schema.ResourceData,
	m any) (b bool, e error) {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamServicesRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
