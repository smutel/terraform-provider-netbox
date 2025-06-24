package ipam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

// Type of vm interface in Netbox
const vMInterfaceType string = "virtualization.vminterface"
const deviceInterfaceType string = "dcim.interface"
const fhrpgroupType string = "ipam.fhrpgroup"

func getNewAvailableIPForIPRange(ctx context.Context,
	client *netbox.APIClient, id int32) (*netbox.AvailableIP,
	diag.Diagnostics) {

	list, response, err := client.IpamAPI.IpamIpRangesAvailableIpsList(
		ctx, id).Execute()

	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &list[0], nil
}

func getNewAvailableIPForPrefix(ctx context.Context,
	client *netbox.APIClient, id int32) (*netbox.AvailableIP,
	diag.Diagnostics) {

	list, response, err := client.IpamAPI.IpamPrefixesAvailableIpsList(
		ctx, id).Execute()

	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &list[0], nil
}

func getNewAvailablePrefix(ctx context.Context,
	client *netbox.APIClient, id int32) (*netbox.AvailablePrefix,
	diag.Diagnostics) {

	resources, response, err :=
		client.IpamAPI.IpamPrefixesAvailablePrefixesList(
			ctx, id).Execute()

	if err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	return &resources[0], nil
}
