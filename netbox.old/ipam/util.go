package ipam

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	netbox "github.com/netbox-community/go-netbox/v4"
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
			if i.GetVirtualMachine != nil {
				return i.GetVirtualMachine().Id, nil
			}
		}
	}

	return 0, util.GenerateErrorMessage(nil, fmt.Errorf("Virtual machine not found"))
}

func isprimary(ctx context.Context, client *netbox.APIClient, objectID int64, ipID int32, ip4 bool) (bool, diag.Diagnostics) {

	if objectID == 0 {
		return false, nil
	}

	var vm *netbox.PaginatedVirtualMachineWithConfigContextList
	var response *http.Response
	var err error
	objectIDArray := []int32{int32(objectID)}
	if ip4 {
		vm, response, err = client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).PrimaryIp4Id(objectIDArray).Execute()
	} else {
		vm, response, err = client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).PrimaryIp6Id(objectIDArray).Execute()
	}

	if err != nil {
		return false, util.GenerateErrorMessage(response, err)
	}

	if vm.GetCount() >= 1 {
		return true, nil
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

	resource := netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	if primary {
		resource.SetPrimaryIp4(ipid)
	} else {
		resource.SetPrimaryIp4Nil()
	}

	if _, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, vmid).WritableVirtualMachineWithConfigContextRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}
