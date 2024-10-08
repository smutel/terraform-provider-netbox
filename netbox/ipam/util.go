package ipam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

// Type of vm interface in Netbox
const vMInterfaceType string = "virtualization.vminterface"
const deviceInterfaceType string = "dcim.interface"
const fhrpgroupType string = "ipam.fhrpgroup"

func getNewAvailableIPForIPRange(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.AvailableIP, diag.Diagnostics) {
	list, response, err := client.IpamAPI.IpamIpRangesAvailableIpsList(ctx, id).Execute()

	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &list[0], nil
}

func getNewAvailableIPForPrefix(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.AvailableIP, diag.Diagnostics) {

	list, response, err := client.IpamAPI.IpamPrefixesAvailableIpsList(ctx, id).Execute()

	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &list[0], nil
}

func getNewAvailablePrefix(ctx context.Context, client *netbox.APIClient, id int32) (*netbox.AvailablePrefix, diag.Diagnostics) {
	resources, response, err := client.IpamAPI.IpamPrefixesAvailablePrefixesList(ctx, id).Execute()
	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &resources[0], nil
}

func getVMIDForInterface(ctx context.Context, client *netbox.APIClient, objectID int32) (int32, diag.Diagnostics) {

	requestObjectID := []int32{objectID}

	interfaces, response, err := client.VirtualizationAPI.VirtualizationInterfacesList(ctx).Id(requestObjectID).Execute()

	if err != nil {
		return 0, util.GenerateErrorMessage(response, err)
	}

	for _, i := range interfaces.Results {
		if i.GetId() == objectID {
			return i.GetVirtualMachine().Id, nil
		}
	}

	return 0, util.GenerateErrorMessage(nil, fmt.Errorf("Virtual machine not found"))
}

func isprimary(ctx context.Context, client *netbox.APIClient, objectID int64, ipID int32, ip4 bool) (bool, diag.Diagnostics) {
	fmt.Println("TATA")
	if objectID == 0 {
		fmt.Println("TITI")
		return false, nil
	}

	resource, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Id([]int32{int32(objectID)}).Execute()

	if err != nil {
		return false, util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() > 0 {
		fmt.Println("TOTO")
		r := resource.Results[0]
		if ip4 && r.GetPrimaryIp4().Id == ipID {
			return true, nil
		}

		if !ip4 && r.GetPrimaryIp6().Id == ipID {
			return true, nil
		}
	}

	return false, nil
}

func setPrimaryIP(ctx context.Context, client *netbox.APIClient, addressID int32, objectID int32, objectType string, primary bool) diag.Diagnostics {

	switch objectType {
	case vMInterfaceType:
		vmID, err := getVMIDForInterface(ctx, client, objectID)
		if err != nil {
			return err
		}
		err = updatePrimaryStatus(ctx, client, vmID, addressID, primary)
		if err != nil {
			return err
		}
		return nil
	case deviceInterfaceType:
		return util.GenerateErrorMessage(nil, fmt.Errorf("this provider does not support the primary_ip4 attribute for '%s'", deviceInterfaceType))
	case fhrpgroupType:
		return util.GenerateErrorMessage(nil, fmt.Errorf("netbox does not support the primary_ip4 attribute for '%s'", fhrpgroupType))
	default:
		return util.GenerateErrorMessage(nil, fmt.Errorf("unknown object type '%s'", objectType))
	}
}

func updatePrimaryStatus(ctx context.Context, client *netbox.APIClient, vmid int32, ipid int32, primary bool) diag.Diagnostics {

	oldResource, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil, int32(vmid)).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	newResource := netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	newResource.SetName(oldResource.GetName())
	if oldResource.GetCluster().Id != 0 {
		b, errDiag := brief.GetBriefClusterRequestFromID(client, ctx, oldResource.GetCluster().Id)
		if errDiag != nil {
			return errDiag
		}
		newResource.SetCluster(*b)
	}
	if oldResource.GetSite().Id != 0 {
		b, errDiag := brief.GetBriefSiteRequestFromID(client, ctx, oldResource.GetSite().Id)
		if errDiag != nil {
			return errDiag
		}
		newResource.SetSite(*b)
	}

	if primary {
		b, err := brief.GetBriefIPAdressRequestFromID(client, ctx, ipid)
		if err != nil {
			return err
		}
		newResource.SetPrimaryIp4(*b)
	} else {
		newResource.SetPrimaryIp4Nil()
	}

	if _, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, vmid).WritableVirtualMachineWithConfigContextRequest(*newResource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}
