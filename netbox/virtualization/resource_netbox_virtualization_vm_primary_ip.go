package virtualization

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
)

func ResourceNetboxVirtualizationVMPrimaryIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a primary ip assignment for a VM (virtualization module) resource within Netbox.",
		CreateContext: resourceNetboxVirtualizationVMPrimaryIPCreate,
		ReadContext:   resourceNetboxVirtualizationVMPrimaryIPRead,
		UpdateContext: resourceNetboxVirtualizationVMPrimaryIPUpdate,
		DeleteContext: resourceNetboxVirtualizationVMPrimaryIPDelete,
		Exists:        resourceNetboxVirtualizationVMPrimaryIPExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"virtualmachine_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Default:     nil,
				Description: "ID of the virtual machine.",
			},
			"primary_ip4_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the primary IPv4 address.",
			},
			"primary_ip6_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the primary IPv4 address.",
			},
		},
	}
}

func resourceNetboxVirtualizationVMPrimaryIPCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return diag.Errorf("virtual machine with ID %d does not exist", d.Get("virtualmachine_id"))
	}

	dropFields := []string{
		"created",
		"last_updated",
		"name",
		"cluster",
		"tags",
	}
	emptyFields := make(map[string]interface{})

	vmID := int64(d.Get("virtualmachine_id").(int))

	newResource := &models.WritableVirtualMachineWithConfigContext{}
	if primaryIP4ID := int64(d.Get("primary_ip4_id").(int)); primaryIP4ID != 0 {
		newResource.PrimaryIp4 = &primaryIP4ID
	}
	if primaryIP6ID := int64(d.Get("primary_ip6_id").(int)); primaryIP6ID != 0 {
		newResource.PrimaryIp6 = &primaryIP6ID
	}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(newResource).WithID(vmID)

	resourceCreated, err := client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {

			// Setting this is only needed for imported resources
			if err := d.Set("virtualmachine_id", resource.ID); err != nil {
				return diag.FromErr(err)
			}

			var primaryIP4ID *int64
			primaryIP4ID = nil
			if resource.PrimaryIp4 != nil {
				primaryIP4ID = &resource.PrimaryIp4.ID
			}
			if err := d.Set("primary_ip4_id", primaryIP4ID); err != nil {
				return diag.FromErr(err)
			}

			var primaryIP6ID *int64
			primaryIP6ID = nil
			if resource.PrimaryIp6 != nil {
				primaryIP6ID = &resource.PrimaryIp6.ID
			}
			if err := d.Set("primary_ip6_id", primaryIP6ID); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return diag.Errorf("virtual machine with ID %d does not exist", d.Get("virtualmachine_id"))
	}

	dropFields := []string{
		"created",
		"last_updated",
		"name",
		"cluster",
		"tags",
	}
	emptyFields := make(map[string]interface{})

	params := &models.WritableVirtualMachineWithConfigContext{}

	if d.HasChange("primary_ip4_id") {
		primaryIP4ID := int64(d.Get("primary_ip4_id").(int))
		params.PrimaryIp4 = &primaryIP4ID
		if primaryIP4ID == 0 {
			emptyFields["primary_ip4"] = nil
		}
	}
	if d.HasChange("primary_ip6_id") {
		primaryIP6ID := int64(d.Get("primary_ip6_id").(int))
		params.PrimaryIp6 = &primaryIP6ID
		if primaryIP6ID == 0 {
			emptyFields["primary_ip6"] = nil
		}
	}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	dropFields := []string{
		"created",
		"last_updated",
		"name",
		"cluster",
		"tags",
	}
	emptyFields := map[string]interface{}{
		"primary_ip4": nil,
		"primary_ip6": nil,
	}

	params := &models.WritableVirtualMachineWithConfigContext{}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := int64(d.Get("virtualmachine_id").(int))
	resourceIDString := strconv.FormatInt(resourceID, 10)
	if resourceIDString == "0" {
		resourceIDString = d.Id()
	}
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceIDString)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == resourceIDString {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
