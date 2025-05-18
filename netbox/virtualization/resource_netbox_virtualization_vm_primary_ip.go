package virtualization

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxVirtualizationVMPrimaryIP() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a primary ip assignment for a VM resource " +
			"within Netbox.",
		CreateContext: resourceNetboxVirtualizationVMPrimaryIPCreate,
		ReadContext:   resourceNetboxVirtualizationVMPrimaryIPRead,
		UpdateContext: resourceNetboxVirtualizationVMPrimaryIPUpdate,
		DeleteContext: resourceNetboxVirtualizationVMPrimaryIPDelete,
		Exists:        resourceNetboxVirtualizationVMPrimaryIPExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this primary ip assignment.",
			},
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

func resourceNetboxVirtualizationVMPrimaryIPCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return util.GenerateErrorMessage(nil,
			fmt.Errorf("virtual machine with "+
				"ID %d does not exist", d.Get("virtualmachine_id")))
	}

	vmID := d.Get("virtualmachine_id").(int)
	vmID32, err := safecast.ToInt32(vmID)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	oldResource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil,
			vmID32).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	newResource :=
		netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	newResource.SetName(oldResource.GetName())
	if oldResource.GetCluster().Id != 0 {
		b, errDiag := brief.GetBriefClusterRequestFromID(ctx, client,
			int(oldResource.GetCluster().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetCluster(*b)
	}
	if oldResource.GetSite().Id != 0 {
		b, errDiag := brief.GetBriefSiteRequestFromID(ctx, client,
			int(oldResource.GetSite().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetSite(*b)
	}
	if primaryIP4ID := d.Get("primary_ip4_id").(int); primaryIP4ID != 0 {
		b, err := brief.GetBriefIPAdressRequestFromID(ctx, client, primaryIP4ID)
		if err != nil {
			return err
		}
		newResource.SetPrimaryIp4(*b)
	}

	if primaryIP6ID := d.Get("primary_ip6_id").(int); primaryIP6ID != 0 {
		b, err := brief.GetBriefIPAdressRequestFromID(ctx, client, primaryIP6ID)
		if err != nil {
			return err
		}
		newResource.SetPrimaryIp6(*b)
	}

	_, response, errDiag :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx,
			vmID32).WritableVirtualMachineWithConfigContextRequest(
			*newResource).Execute()
	if response.StatusCode != util.Const201 && errDiag != nil {
		return util.GenerateErrorMessage(response, errDiag)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type",
		util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	// Setting this is only needed for imported resources
	if err := d.Set("virtualmachine_id", resource.GetId()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err := d.Set("primary_ip4_id", resource.GetPrimaryIp4().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err := d.Set("primary_ip6_id", resource.GetPrimaryIp6().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return util.GenerateErrorMessage(nil,
			fmt.Errorf("virtual machine with "+
				"ID %d does not exist", d.Get("virtualmachine_id")))
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	oldResource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil,
			int32(resourceID)).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	newResource :=
		netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	newResource.SetName(oldResource.GetName())
	if oldResource.GetCluster().Id != 0 {
		b, errDiag := brief.GetBriefClusterRequestFromID(ctx, client,
			int(oldResource.GetCluster().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetCluster(*b)
	}
	if oldResource.GetSite().Id != 0 {
		b, errDiag := brief.GetBriefSiteRequestFromID(ctx, client,
			int(oldResource.GetSite().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetSite(*b)
	}

	if d.HasChange("primary_ip4_id") {
		if primaryIP4ID := d.Get("primary_ip4_id").(int); primaryIP4ID != 0 {
			b, err := brief.GetBriefIPAdressRequestFromID(ctx, client,
				primaryIP4ID)
			if err != nil {
				return err
			}
			newResource.SetPrimaryIp4(*b)
		} else {
			newResource.SetPrimaryIp4Nil()
		}
	}

	if d.HasChange("primary_ip6_id") {
		if primaryIP6ID := d.Get("primary_ip6_id").(int); primaryIP6ID != 0 {
			b, err := brief.GetBriefIPAdressRequestFromID(ctx, client,
				primaryIP6ID)
			if err != nil {
				return err
			}
			newResource.SetPrimaryIp6(*b)
		} else {
			newResource.SetPrimaryIp6Nil()
		}
	}

	if _, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx,
			int32(resourceID)).WritableVirtualMachineWithConfigContextRequest(
			*newResource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxVirtualizationVMPrimaryIPRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMPrimaryIPDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMPrimaryIPExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	oldResource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil,
			int32(resourceID)).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	newResource :=
		netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	newResource.SetName(oldResource.GetName())
	if oldResource.GetCluster().Id != 0 {
		b, errDiag := brief.GetBriefClusterRequestFromID(ctx, client,
			int(oldResource.GetCluster().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetCluster(*b)
	}
	if oldResource.GetSite().Id != 0 {
		b, errDiag := brief.GetBriefSiteRequestFromID(ctx, client,
			int(oldResource.GetSite().Id))
		if errDiag != nil {
			return errDiag
		}
		newResource.SetSite(*b)
	}

	if primaryIP4ID := d.Get("primary_ip4_id").(int); primaryIP4ID != 0 {
		newResource.SetPrimaryIp4Nil()
	}

	if primaryIP6ID := d.Get("primary_ip6_id").(int); primaryIP6ID != 0 {
		newResource.SetPrimaryIp6Nil()
	}

	if _, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx,
			int32(resourceID)).WritableVirtualMachineWithConfigContextRequest(
			*newResource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxVirtualizationVMPrimaryIPExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID := d.Get("virtualmachine_id").(int)
	resourceID32, err := safecast.ToInt32(resourceID)
	if err != nil {
		return false, err
	}

	_, http, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil,
			resourceID32).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
