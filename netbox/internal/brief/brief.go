package brief

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func GetBriefManufacturerRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefManufacturerRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimManufacturersRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefManufacturerRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefDeviceRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefDeviceRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimDevicesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefDeviceRequest()
	m.SetName(resource.GetName())

	return m, nil
}

func GetBriefDeviceRoleRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefDeviceRoleRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimDeviceRolesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefDeviceRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefPlatformRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefPlatformRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimPlatformsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefPlatformRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefConfigTemplateRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefConfigTemplateRequest, diag.Diagnostics) {
	resource, response, err := client.ExtrasAPI.ExtrasConfigTemplatesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefConfigTemplateRequest(resource.GetName())

	return m, nil
}

func GetBriefSiteGroupRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefSiteGroupRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimSiteGroupsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefSiteGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefRegionRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefRegionRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimRegionsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRegionRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefTenantRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefTenantRequest, diag.Diagnostics) {
	resource, response, err := client.TenancyAPI.TenancyTenantsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefTenantRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefTenantGroupRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefTenantGroupRequest, diag.Diagnostics) {
	resource, response, err := client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefTenantGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefRIRRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefRIRRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamRirsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRIRRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefIPAdressRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefIPAddressRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamIpAddressesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefIPAddressRequest(resource.GetAddress())

	return m, nil
}

func GetBriefVRFRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefVRFRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamVrfsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVRFRequest(resource.GetName())

	return m, nil
}

func GetBriefVLANRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefVLANRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamVlansRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVLANRequest(resource.GetVid(), resource.GetName())

	return m, nil
}

func GetBriefVlanRoleRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefRoleRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamRolesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefVlanGroupRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefVLANGroupRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamVlanGroupsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVLANGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefVlanRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefVLANRequest, diag.Diagnostics) {
	resource, response, err := client.IpamAPI.IpamVlansRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVLANRequest(resource.GetVid(), resource.GetName())

	return m, nil
}

func GetBriefSiteRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefSiteRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimSitesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefSiteRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefLocationRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefLocationRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimLocationsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefLocationRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefRackRoleRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefRackRoleRequest, diag.Diagnostics) {
	resource, response, err := client.DcimAPI.DcimRackRolesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRackRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefContactGroupRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefContactGroupRequest, diag.Diagnostics) {
	resource, response, err := client.TenancyAPI.TenancyContactGroupsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefContactRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefContactRequest, diag.Diagnostics) {
	resource, response, err := client.TenancyAPI.TenancyContactsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactRequest(resource.GetName())

	return m, nil
}

func GetBriefContactRoleRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefContactRoleRequest, diag.Diagnostics) {
	resource, response, err := client.TenancyAPI.TenancyContactRolesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefClusterRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefClusterRequest, diag.Diagnostics) {
	resource, response, err := client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterRequest(resource.GetName())

	return m, nil
}

func GetBriefClusterTypeRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefClusterTypeRequest, diag.Diagnostics) {
	resource, response, err := client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterTypeRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefClusterGroupRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefClusterGroupRequest, diag.Diagnostics) {
	resource, response, err := client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefVirtualMachineRequestFromID(client *netbox.APIClient, ctx context.Context, id int32) (*netbox.BriefVirtualMachineRequest, diag.Diagnostics) {
	resource, response, err := client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, id).Execute()

	if response.StatusCode == 404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVirtualMachineRequest(resource.GetName())

	return m, nil
}
