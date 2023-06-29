package virtualization

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

// Type of vm interface in Netbox
const vMInterfaceType string = "virtualization.vminterface"

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
	client := m.(*netboxclient.NetBoxAPI)

	dropFields := []string{
		"created",
		"last_updated",
	}
	emptyFields := make(map[string]interface{})

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	mode := d.Get("mode").(string)
	mtu := int64(d.Get("mtu").(int))
	name := d.Get("name").(string)
	taggedVlans, err := util.ExpandToInt64Slice(d.Get("tagged_vlans").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	tags := d.Get("tag").(*schema.Set).List()
	untaggedVlan := int64(d.Get("untagged_vlan").(int))
	virtualmachineID := int64(d.Get("virtualmachine_id").(int))

	newResource := &models.WritableVMInterface{
		CustomFields:   &customFields,
		Description:    description,
		Enabled:        enabled,
		Mode:           mode,
		Name:           &name,
		TaggedVlans:    taggedVlans,
		Tags:           tag.ConvertTagsToNestedTags(tags),
		VirtualMachine: &virtualmachineID,
	}

	if !enabled {
		emptyFields["enabled"] = false
	}

	if macAddress := d.Get("mac_address").(string); macAddress != "" {
		newResource.MacAddress = &macAddress
	}

	if bridgeID := int64(d.Get("bridge_id").(int)); bridgeID != 0 {
		newResource.Bridge = &bridgeID
	}

	if parentID := int64(d.Get("parent_id").(int)); parentID != 0 {
		newResource.Parent = &parentID
	}

	if vrfID := int64(d.Get("vrf_id").(int)); vrfID != 0 {
		newResource.Vrf = &vrfID
	}

	if mtu != 0 {
		newResource.Mtu = &mtu
	}

	if untaggedVlan != 0 {
		newResource.UntaggedVlan = &untaggedVlan
	}

	resource := virtualization.NewVirtualizationInterfacesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationInterfacesCreate(
		resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationInterfaceRead(ctx, d, m)
}

func resourceNetboxVirtualizationInterfaceRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationInterfacesList(
		params, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}
			var bridgeID *int64
			bridgeID = nil
			if resource.Bridge != nil {
				bridgeID = &resource.Bridge.ID
			}
			if err = d.Set("bridge_id", bridgeID); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("count_fhrp_groups", resource.CountFhrpGroups); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("count_ipaddresses", resource.CountIpaddresses); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("created", resource.Created.String()); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("description", resource.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("enabled", resource.Enabled); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("mac_address", resource.MacAddress); err != nil {
				return diag.FromErr(err)
			}

			var mode *string
			mode = nil
			if resource.Mode != nil {
				mode = resource.Mode.Value
			}
			if err = d.Set("mode", mode); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("mtu", resource.Mtu); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}
			var parentID *int64
			parentID = nil
			if resource.Parent != nil {
				parentID = &resource.Parent.ID
			}
			if err = d.Set("parent_id", parentID); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tagged_vlans", util.ConvertNestedVlansToVlans(resource.TaggedVlans)); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", tag.ConvertNestedTagsToTags(
				resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("untagged_vlan", util.GetNestedVlanID(resource.UntaggedVlan)); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("virtualmachine_id",
				resource.VirtualMachine.ID); err != nil {
				return diag.FromErr(err)
			}
			var vrfID *int64
			vrfID = nil
			if resource.Vrf != nil {
				vrfID = &resource.Vrf.ID
			}
			if err = d.Set("vrf_id", vrfID); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("type", vMInterfaceType); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationInterfaceUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVMInterface{}
	dropFields := []string{
		"created",
		"last_updated",
	}
	emptyFields := make(map[string]interface{})

	if d.HasChange("bridge_id") {
		bridgeID := int64(d.Get("bridge_id").(int))
		if bridgeID != 0 {
			params.Bridge = &bridgeID
		} else {
			emptyFields["bridge"] = nil
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		if description != "" {
			params.Description = description
		} else {
			emptyFields["description"] = ""
		}
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			params.Enabled = enabled
		} else {
			emptyFields["enabled"] = false
		}
	}

	if d.HasChange("mac_address") {
		macAddress := d.Get("mac_address").(string)
		if macAddress != "" {
			params.MacAddress = &macAddress
		} else {
			emptyFields["mac_address"] = nil
		}
	}

	if d.HasChange("mode") {
		mode := d.Get("mode").(string)
		if mode != "" {
			params.Mode = mode
		} else {
			emptyFields["mode"] = nil
		}
	}

	if d.HasChange("mtu") {
		mtu := int64(d.Get("mtu").(int))
		if mtu != 0 {
			params.Mtu = &mtu
		} else {
			emptyFields["mtu"] = nil
		}
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	} else {
		dropFields = append(dropFields, "name")
	}

	if d.HasChange("parent_id") {
		parentID := int64(d.Get("parent_id").(int))
		if parentID != 0 {
			params.Parent = &parentID
		} else {
			emptyFields["parent"] = nil
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	} else {
		dropFields = append(dropFields, "tags")
	}

	if d.HasChange("tagged_vlans") {
		taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
		tvlans, err := util.ExpandToInt64Slice(taggedVlans)
		if err != nil {
			return diag.FromErr(err)
		}
		params.TaggedVlans = tvlans
	} else {
		dropFields = append(dropFields, "tagged_vlans")
	}

	if d.HasChange("untagged_vlan") {
		untaggedVlan := int64(d.Get("untagged_vlan").(int))
		params.UntaggedVlan = &untaggedVlan
		if untaggedVlan == 0 {
			emptyFields["untagged_vlan"] = nil
		}
	}

	if d.HasChange("virtual_machine_id") {
		virtualMachineID := int64(d.Get("virtual_machine_id").(int))
		if virtualMachineID != 0 {
			params.VirtualMachine = &virtualMachineID
		}
	} else {
		dropFields = append(dropFields, "virtual_machine")
	}

	if d.HasChange("vrf_id") {
		vrfID := int64(d.Get("vrf_id").(int))
		if vrfID != 0 {
			params.Vrf = &vrfID
		} else {
			emptyFields["vrf"] = nil
		}
	}

	resource := virtualization.NewVirtualizationInterfacesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationInterfacesPartialUpdate(
		resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxVirtualizationInterfaceRead(ctx, d, m)
}

func resourceNetboxVirtualizationInterfaceDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationInterfaceExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	p := virtualization.NewVirtualizationInterfacesDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationInterfacesDelete(
		p, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationInterfaceExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationInterfacesList(
		params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
