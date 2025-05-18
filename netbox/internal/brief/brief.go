package brief

import (
	"context"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func GetBriefManufacturerRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefManufacturerRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimManufacturersRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefManufacturerRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefDeviceRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefDeviceRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimDevicesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefDeviceRequest()
	m.SetName(resource.GetName())

	return m, nil
}

func GetBriefDeviceRoleRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefDeviceRoleRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimDeviceRolesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefDeviceRoleRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefPlatformRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefPlatformRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimPlatformsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefPlatformRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefConfigTemplateRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefConfigTemplateRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.ExtrasAPI.ExtrasConfigTemplatesRetrieve(
		ctx, id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefConfigTemplateRequest(resource.GetName())

	return m, nil
}

func GetBriefSiteGroupRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefSiteGroupRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimSiteGroupsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefSiteGroupRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefRegionRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefRegionRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimRegionsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRegionRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefTenantRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefTenantRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.TenancyAPI.TenancyTenantsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefTenantRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefTenantGroupRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefTenantGroupRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.TenancyAPI.TenancyTenantGroupsRetrieve(
		ctx, id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefTenantGroupRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefRIRRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefRIRRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamRirsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRIRRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefIPAdressRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefIPAddressRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamIpAddressesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefIPAddressRequest(resource.GetAddress())

	return m, nil
}

func GetBriefVRFRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefVRFRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamVrfsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVRFRequest(resource.GetName())

	return m, nil
}

func GetBriefVLANRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefVLANRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamVlansRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVLANRequest(resource.GetVid(), resource.GetName())

	return m, nil
}

func GetBriefVlanRoleRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefRoleRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamRolesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefVlanGroupRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefVLANGroupRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.IpamAPI.IpamVlanGroupsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVLANGroupRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefSiteRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefSiteRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimSitesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefSiteRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefLocationRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefLocationRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimLocationsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefLocationRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefRackRoleRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefRackRoleRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.DcimAPI.DcimRackRolesRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefRackRoleRequest(resource.GetName(), resource.GetSlug())

	return m, nil
}

func GetBriefContactGroupRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefContactGroupRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.TenancyAPI.TenancyContactGroupsRetrieve(
		ctx, id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactGroupRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefContactRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefContactRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err := client.TenancyAPI.TenancyContactsRetrieve(ctx,
		id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactRequest(resource.GetName())

	return m, nil
}

func GetBriefContactRoleRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefContactRoleRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err :=
		client.TenancyAPI.TenancyContactRolesRetrieve(ctx, id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefContactRoleRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefClusterRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefClusterRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationClustersRetrieve(
			ctx, id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterRequest(resource.GetName())

	return m, nil
}

func GetBriefClusterTypeRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefClusterTypeRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx,
			id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterTypeRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefClusterGroupRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefClusterGroupRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx,
			id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefClusterGroupRequest(resource.GetName(),
		resource.GetSlug())

	return m, nil
}

func GetBriefVirtualMachineRequestFromID(ctx context.Context,
	client *netbox.APIClient, id int) (
	*netbox.BriefVirtualMachineRequest, diag.Diagnostics) {

	id32, err := safecast.ToInt32(id)
	if err != nil {
		return nil, util.GenerateErrorMessage(nil, err)
	}

	resource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx,
			id32).Execute()

	if response.StatusCode == util.Const404 || err != nil {
		return nil, util.GenerateErrorMessage(response, err)
	}

	m := netbox.NewBriefVirtualMachineRequest(resource.GetName())

	return m, nil
}
