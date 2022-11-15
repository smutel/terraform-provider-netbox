package netbox

import (
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func getNewAvailableIPForIPRange(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamIPRangesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamIPRangesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailableIPForPrefix(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamPrefixesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamPrefixesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailablePrefix(client *netboxclient.NetBoxAPI, id int64, length int64) (*models.Prefix, error) {
	params := ipam.NewIpamPrefixesAvailablePrefixesCreateParams().WithID(id)
	params.Data = []*models.PrefixLength{
		{PrefixLength: &length},
	}
	list, err := client.Ipam.IpamPrefixesAvailablePrefixesCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}
