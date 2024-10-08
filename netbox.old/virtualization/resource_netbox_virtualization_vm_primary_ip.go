package virtualization

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
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
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return util.GenerateErrorMessage(nil, fmt.Errorf("virtual machine with ID %d does not exist", d.Get("virtualmachine_id")))
	}

	vmID := int32(d.Get("virtualmachine_id").(int))

	newResource := netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	if primaryIP4ID := int32(d.Get("primary_ip4_id").(int)); primaryIP4ID != 0 {
		newResource.SetPrimaryIp4(primaryIP4ID)
	}

	if primaryIP6ID := int32(d.Get("primary_ip6_id").(int)); primaryIP6ID != 0 {
		newResource.SetPrimaryIp6(primaryIP6ID)
	}

	resourceCreated, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, vmID).WritableVirtualMachineWithConfigContextRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// Setting this is only needed for imported resources
	if err := d.Set("virtualmachine_id", resource.GetId()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err := d.Set("primary_ip4_id", resource.GetPrimaryIp4().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err := d.Set("primary_ip6_id", resource.GetPrimaryIp6().Id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return util.GenerateErrorMessage(nil, fmt.Errorf("virtual machine with ID %d does not exist", d.Get("virtualmachine_id")))
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()

	if d.HasChange("primary_ip4_id") {
		if primaryIP4ID := int32(d.Get("primary_ip4_id").(int)); primaryIP4ID != 0 {
			resource.SetPrimaryIp4(primaryIP4ID)
		} else {
			resource.SetPrimaryIp4Nil()
		}
	}

	if d.HasChange("primary_ip6_id") {
		if primaryIP6ID := int32(d.Get("primary_ip6_id").(int)); primaryIP6ID != 0 {
			resource.SetPrimaryIp6(primaryIP6ID)
		} else {
			resource.SetPrimaryIp6Nil()
		}
	}

	if _, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, int32(resourceID)).WritableVirtualMachineWithConfigContextRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
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

	resource := netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()

	if primaryIP4ID := int32(d.Get("primary_ip4_id").(int)); primaryIP4ID != 0 {
		resource.SetPrimaryIp4Nil()
	}

	if primaryIP6ID := int32(d.Get("primary_ip6_id").(int)); primaryIP6ID != 0 {
		resource.SetPrimaryIp6Nil()
	}

	if _, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, int32(resourceID)).WritableVirtualMachineWithConfigContextRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID := d.Get("virtualmachine_id").(int)
	if resourceID == 0 {
		var err error
		resourceID, err = strconv.Atoi(d.Id())
		if err != nil {
			return false, err
		}
	}

	_, http, err := client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
