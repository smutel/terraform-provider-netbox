package virtualization

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

// Type of vm interface in Netbox
const vmIntefaceType string = "virtualization.vminterface"

func ResourceNetboxVirtualizationInterface() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an interface (virtualization module) resource within Netbox.",
		CreateContext: resourceNetboxVirtualizationInterfaceCreate,
		ReadContext:   resourceNetboxVirtualizationInterfaceRead,
		UpdateContext: resourceNetboxVirtualizationInterfaceUpdate,
		DeleteContext: resourceNetboxVirtualizationInterfaceDelete,
		Exists:        resourceNetboxVirtualizationInterfaceExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bridge_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the bridge interface where this interface (virtualization module) is attached to.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this interface (virtualization module).",
			},
			"count_fhrp_groups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of fhrp groups attached to this interface (virtualization module) is attached to.",
			},
			"count_ipaddresses": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of ip addresses attached to this interface (virtualization module) is attached to.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this resource was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "Description for this interface (virtualization module).",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "true or false (true by default)",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this resource was last updated.",
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^([A-Z0-9]{2}:){5}[A-Z0-9]{2}$"),
					"Must be like AA:AA:AA:AA:AA"),
				Description: "Mac address for this interface (virtualization module)",
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"access", "tagged",
					"tagged-all"}, false),
				Description: "The mode among access, tagged, tagged-all.",
			},
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65536),
				Description:  "The MTU between 1 and 65536 for this interface (virtualization module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "Description for this interface (virtualization module)",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the parent interface where this interface (virtualization module) is attached to.",
			},
			"tagged_vlans": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "List of vlan id tagged for this interface (virtualization module)",
			},
			"tag": &tag.TagSchema,
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of interface among virtualization.vminterface for VM or dcim.interface for device",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this interface (virtualization module).",
			},
			"untagged_vlan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Vlan ID untagged for this interface (virtualization module).",
			},
			"virtualmachine_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the VM where this interface (virtualization module) is attached to.",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the VRF where this interface (virtualization module) is attached to.",
			},
		},
	}
}

func resourceNetboxVirtualizationInterfaceCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	mode := d.Get("mode").(string)
	mtu := int32(d.Get("mtu").(int))
	name := d.Get("name").(string)
	taggedVlans, err := util.ExpandToInt32Slice(d.Get("tagged_vlans").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	tags := d.Get("tag").(*schema.Set).List()
	untaggedVlan := int32(d.Get("untagged_vlan").(int))
	virtualmachineID := int32(d.Get("virtualmachine_id").(int))

	newResource := netbox.NewWritableVMInterfaceRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetEnabled(enabled)
	newResource.SetName(name)
	newResource.SetTaggedVlans(taggedVlans)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	newResource.SetVirtualMachine(virtualmachineID)

	modeObj, err := netbox.NewPatchedWritableInterfaceRequestModeFromValue(mode)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetMode(*modeObj)

	if macAddress := d.Get("mac_address").(string); macAddress != "" {
		newResource.SetMacAddress(macAddress)
	}

	if bridgeID := int32(d.Get("bridge_id").(int)); bridgeID != 0 {
		newResource.SetBridge(bridgeID)
	}

	if parentID := int32(d.Get("parent_id").(int)); parentID != 0 {
		newResource.SetParent(parentID)
	}

	if vrfID := int32(d.Get("vrf_id").(int)); vrfID != 0 {
		newResource.SetVrf(vrfID)
	}

	if mtu != 0 {
		newResource.SetMtu(mtu)
	}

	if untaggedVlan != 0 {
		newResource.SetUntaggedVlan(untaggedVlan)
	}

	resourceCreated, response, err := client.VirtualizationAPI.VirtualizationInterfacesCreate(ctx).WritableVMInterfaceRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxVirtualizationInterfaceRead(ctx, d, m)
}

func resourceNetboxVirtualizationInterfaceRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, int32(resourceID)).Execute()

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

	if err = d.Set("bridge_id", resource.Bridge.Get().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("count_fhrp_groups", resource.GetCountFhrpGroups()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("count_ipaddresses", resource.GetCountIpaddresses()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated()); err != nil {
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

	if err = d.Set("enabled", resource.GetEnabled()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("mac_address", resource.GetMacAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("mode", resource.GetMode().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("mtu", resource.GetMtu()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("parent_id", resource.GetParent().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tagged_vlans", resource.GetTaggedVlans()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("untagged_vlan", resource.GetUntaggedVlan().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_id", resource.GetVirtualMachine().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vrf_id", resource.GetVrf().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("type", vmIntefaceType); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxVirtualizationInterfaceUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableVMInterfaceRequestWithDefaults()

	if d.HasChange("bridge_id") {
		if bridgeID := int32(d.Get("bridge_id").(int)); bridgeID != 0 {
			resource.SetBridge(bridgeID)
		} else {
			resource.SetBridgeNil()
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("enabled") {
		resource.SetEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("mac_address") {
		macAddress := d.Get("mac_address").(string)
		if macAddress != "" {
			resource.SetMacAddress(macAddress)
		} else {
			resource.SetMacAddressNil()
		}
	}

	if d.HasChange("mode") {
		modeObj, err := netbox.NewPatchedWritableInterfaceRequestModeFromValue(d.Get("mode").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetMode(*modeObj)
	}

	if d.HasChange("mtu") {
		if mtu := int32(d.Get("mtu").(int)); mtu != 0 {
			resource.SetMtu(mtu)
		} else {
			resource.SetMtuNil()
		}
	}

	if d.HasChange("name") {
		resource.SetName(d.Get("name").(string))
	}

	if d.HasChange("parent_id") {
		if parentID := int32(d.Get("parent_id").(int)); parentID != 0 {
			resource.SetParent(parentID)
		} else {
			resource.SetParentNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tagged_vlans") {
		taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
		tvlans, err := util.ExpandToInt32Slice(taggedVlans)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetTaggedVlans(tvlans)
	}

	if d.HasChange("untagged_vlan") {
		if untaggedVlan := int32(d.Get("untagged_vlan").(int)); untaggedVlan != 0 {
			resource.SetUntaggedVlan(untaggedVlan)
		} else {
			resource.SetUntaggedVlanNil()
		}
	}

	if d.HasChange("virtual_machine_id") {
		resource.SetVirtualMachine(int32(d.Get("virtual_machine_id").(int)))
	}

	if d.HasChange("vrf_id") {
		if vrfID := int32(d.Get("vrf_id").(int)); vrfID != 0 {
			resource.SetVrf(vrfID)
		} else {
			resource.SetVrfNil()
		}
	}

	if _, response, err := client.VirtualizationAPI.VirtualizationInterfacesUpdate(ctx, int32(resourceID)).WritableVMInterfaceRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxVirtualizationInterfaceRead(ctx, d, m)
}

func resourceNetboxVirtualizationInterfaceDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationInterfaceExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}

	if response, err := client.VirtualizationAPI.VirtualizationInterfacesDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxVirtualizationInterfaceExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.VirtualizationAPI.VirtualizationInterfacesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
