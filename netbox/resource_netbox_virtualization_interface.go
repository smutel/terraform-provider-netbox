package netbox

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
)

func resourceNetboxVirtualizationInterface() *schema.Resource {
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
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this interface (virtualization module).",
			},
			"custom_field": &customFieldSchema,
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
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^([A-Z0-9]{2}:){5}[A-Z0-9]{2}$"),
					"Must be like AA:AA:AA:AA:AA"),
				ForceNew:    true,
				Description: "Mac address for this interface (virtualization module)",
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"access", "tagged",
					"tagged-all"}, false),
				ForceNew:    true,
				Description: "The mode among access, tagged, tagged-all.",
			},
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65536),
				ForceNew:     true,
				Description:  "The MTU between 1 and 65536 for this interface (virtualization module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "Description for this interface (virtualization module)",
			},
			"tagged_vlans": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "List of vlan id tagged for this interface (virtualization module)",
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing tag.",
						},
						"slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Slug of the existing tag.",
						},
					},
				},
				Description: "Existing tag to associate to this interface (virtualization module).",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of interface among virtualization.vminterface for VM or dcim.interface for device",
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
		},
	}
}

func resourceNetboxVirtualizationInterfaceCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	macAddress := d.Get("mac_address").(string)
	mode := d.Get("mode").(string)
	mtu := int64(d.Get("mtu").(int))
	name := d.Get("name").(string)
	taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
	tags := d.Get("tag").(*schema.Set).List()
	untaggedVlan := int64(d.Get("untagged_vlan").(int))
	virtualmachineID := int64(d.Get("virtualmachine_id").(int))

	newResource := &models.WritableVMInterface{
		CustomFields:   &customFields,
		Description:    description,
		Enabled:        enabled,
		Mode:           mode,
		Name:           &name,
		TaggedVlans:    expandToInt64Slice(taggedVlans),
		Tags:           convertTagsToNestedTags(tags),
		VirtualMachine: &virtualmachineID,
	}

	if macAddress != "" {
		newResource.MacAddress = &macAddress
	}

	if mtu != 0 {
		newResource.Mtu = &mtu
	}

	if untaggedVlan != 0 {
		newResource.UntaggedVlan = &untaggedVlan
	}

	resource := virtualization.NewVirtualizationInterfacesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationInterfacesCreate(
		resource, nil)

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
			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("enabled", resource.Enabled); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("mac_address", resource.MacAddress); err != nil {
				return diag.FromErr(err)
			}

			if resource.Mode == nil {
				if err = d.Set("mode", ""); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("mode", resource.Mode.Value); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("mtu", resource.Mtu); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tagged_vlans", resource.TaggedVlans); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(
				resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("untagged_vlan", resource.UntaggedVlan); err != nil {
				return diag.FromErr(err)
			}

			if resource.VirtualMachine == nil {
				if err = d.Set("virtualmachine_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("virtualmachine_id",
					resource.VirtualMachine.ID); err != nil {
					return diag.FromErr(err)
				}
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

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name
	virtualMachineID := int64(d.Get("virtualmachine_id").(int))
	params.VirtualMachine = &virtualMachineID
	taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
	params.TaggedVlans = expandToInt64Slice(taggedVlans)

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

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		params.Enabled = enabled
	}

	if d.HasChange("mac_address") {
		macAddress := d.Get("mac_address").(string)
		params.MacAddress = &macAddress
	}

	if d.HasChange("mode") {
		mode := d.Get("mode").(string)
		params.Mode = mode
	}

	if d.HasChange("mtu") {
		mtu := int64(d.Get("mtu").(int))
		params.Mtu = &mtu
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("untagged_vlan") {
		untaggedVlan := int64(d.Get("untagged_vlan").(int))
		params.UntaggedVlan = &untaggedVlan
	}

	resource := virtualization.NewVirtualizationInterfacesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationInterfacesPartialUpdate(
		resource, nil)
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
