package netbox

import (
	"fmt"

	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client"
)

const authHeaderName = "Authorization"
const authHeaderFormat = "Token %v"

// Provider exports the actual provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_URL", "127.0.0.1:8000"),
				Description: "URL to reach netbox application.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_TOKEN", ""),
				Description: "Token used for API operations.",
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_SCHEME", "https"),
				Description: "Sheme used to reach netbox application.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netbox_dcim_site":            dataNetboxDcimSite(),
			"netbox_ipam_ip_addresses":    dataNetboxIpamIPAddresses(),
			"netbox_ipam_role":            dataNetboxIpamRole(),
			"netbox_ipam_vlan":            dataNetboxIpamVlan(),
			"netbox_ipam_vlan_group":      dataNetboxIpamVlanGroup(),
			"netbox_tenancy_tenant":       dataNetboxTenancyTenant(),
			"netbox_tenancy_tenant_group": dataNetboxTenancyTenantGroup(),
			"netbox_json_circuits_circuit-terminations_list": dataNetboxJsonCircuitsCircuitTerminationsList(),
			"netbox_json_circuits_circuit-types_list": dataNetboxJsonCircuitsCircuitTypesList(),
			"netbox_json_circuits_circuits_list": dataNetboxJsonCircuitsCircuitsList(),
			"netbox_json_circuits_providers_list": dataNetboxJsonCircuitsProvidersList(),
			"netbox_json_dcim_cables_list": dataNetboxJsonDcimCablesList(),
			"netbox_json_dcim_console-connections_list": dataNetboxJsonDcimConsoleConnectionsList(),
			"netbox_json_dcim_console-port-templates_list": dataNetboxJsonDcimConsolePortTemplatesList(),
			"netbox_json_dcim_console-ports_list": dataNetboxJsonDcimConsolePortsList(),
			"netbox_json_dcim_console-server-port-templates_list": dataNetboxJsonDcimConsoleServerPortTemplatesList(),
			"netbox_json_dcim_console-server-ports_list": dataNetboxJsonDcimConsoleServerPortsList(),
			"netbox_json_dcim_device-bay-templates_list": dataNetboxJsonDcimDeviceBayTemplatesList(),
			"netbox_json_dcim_device-bays_list": dataNetboxJsonDcimDeviceBaysList(),
			"netbox_json_dcim_device-roles_list": dataNetboxJsonDcimDeviceRolesList(),
			"netbox_json_dcim_device-types_list": dataNetboxJsonDcimDeviceTypesList(),
			"netbox_json_dcim_devices_list": dataNetboxJsonDcimDevicesList(),
			"netbox_json_dcim_front-port-templates_list": dataNetboxJsonDcimFrontPortTemplatesList(),
			"netbox_json_dcim_front-ports_list": dataNetboxJsonDcimFrontPortsList(),
			"netbox_json_dcim_interface-connections_list": dataNetboxJsonDcimInterfaceConnectionsList(),
			"netbox_json_dcim_interface-templates_list": dataNetboxJsonDcimInterfaceTemplatesList(),
			"netbox_json_dcim_interfaces_list": dataNetboxJsonDcimInterfacesList(),
			"netbox_json_dcim_inventory-items_list": dataNetboxJsonDcimInventoryItemsList(),
			"netbox_json_dcim_manufacturers_list": dataNetboxJsonDcimManufacturersList(),
			"netbox_json_dcim_platforms_list": dataNetboxJsonDcimPlatformsList(),
			"netbox_json_dcim_power-connections_list": dataNetboxJsonDcimPowerConnectionsList(),
			"netbox_json_dcim_power-feeds_list": dataNetboxJsonDcimPowerFeedsList(),
			"netbox_json_dcim_power-outlet-templates_list": dataNetboxJsonDcimPowerOutletTemplatesList(),
			"netbox_json_dcim_power-outlets_list": dataNetboxJsonDcimPowerOutletsList(),
			"netbox_json_dcim_power-panels_list": dataNetboxJsonDcimPowerPanelsList(),
			"netbox_json_dcim_power-port-templates_list": dataNetboxJsonDcimPowerPortTemplatesList(),
			"netbox_json_dcim_power-ports_list": dataNetboxJsonDcimPowerPortsList(),
			"netbox_json_dcim_rack-groups_list": dataNetboxJsonDcimRackGroupsList(),
			"netbox_json_dcim_rack-reservations_list": dataNetboxJsonDcimRackReservationsList(),
			"netbox_json_dcim_rack-roles_list": dataNetboxJsonDcimRackRolesList(),
			"netbox_json_dcim_racks_list": dataNetboxJsonDcimRacksList(),
			"netbox_json_dcim_rear-port-templates_list": dataNetboxJsonDcimRearPortTemplatesList(),
			"netbox_json_dcim_rear-ports_list": dataNetboxJsonDcimRearPortsList(),
			"netbox_json_dcim_regions_list": dataNetboxJsonDcimRegionsList(),
			"netbox_json_dcim_sites_list": dataNetboxJsonDcimSitesList(),
			"netbox_json_dcim_virtual-chassis_list": dataNetboxJsonDcimVirtualChassisList(),
			"netbox_json_extras_config-contexts_list": dataNetboxJsonExtrasConfigContextsList(),
			"netbox_json_extras_export-templates_list": dataNetboxJsonExtrasExportTemplatesList(),
			"netbox_json_extras_graphs_list": dataNetboxJsonExtrasGraphsList(),
			"netbox_json_extras_image-attachments_list": dataNetboxJsonExtrasImageAttachmentsList(),
			"netbox_json_extras_job-results_list": dataNetboxJsonExtrasJobResultsList(),
			"netbox_json_extras_object-changes_list": dataNetboxJsonExtrasObjectChangesList(),
			"netbox_json_extras_tags_list": dataNetboxJsonExtrasTagsList(),
			"netbox_json_ipam_aggregates_list": dataNetboxJsonIpamAggregatesList(),
			"netbox_json_ipam_ip-addresses_list": dataNetboxJsonIpamIPAddressesList(),
			"netbox_json_ipam_prefixes_list": dataNetboxJsonIpamPrefixesList(),
			"netbox_json_ipam_rirs_list": dataNetboxJsonIpamRirsList(),
			"netbox_json_ipam_roles_list": dataNetboxJsonIpamRolesList(),
			"netbox_json_ipam_services_list": dataNetboxJsonIpamServicesList(),
			"netbox_json_ipam_vlan-groups_list": dataNetboxJsonIpamVlanGroupsList(),
			"netbox_json_ipam_vlans_list": dataNetboxJsonIpamVlansList(),
			"netbox_json_ipam_vrfs_list": dataNetboxJsonIpamVrfsList(),
			"netbox_json_secrets_secret-roles_list": dataNetboxJsonSecretsSecretRolesList(),
			"netbox_json_secrets_secrets_list": dataNetboxJsonSecretsSecretsList(),
			"netbox_json_tenancy_tenant-groups_list": dataNetboxJsonTenancyTenantGroupsList(),
			"netbox_json_tenancy_tenants_list": dataNetboxJsonTenancyTenantsList(),
			"netbox_json_users_groups_list": dataNetboxJsonUsersGroupsList(),
			"netbox_json_users_permissions_list": dataNetboxJsonUsersPermissionsList(),
			"netbox_json_users_users_list": dataNetboxJsonUsersUsersList(),
			"netbox_json_virtualization_cluster-groups_list": dataNetboxJsonVirtualizationClusterGroupsList(),
			"netbox_json_virtualization_cluster-types_list": dataNetboxJsonVirtualizationClusterTypesList(),
			"netbox_json_virtualization_clusters_list": dataNetboxJsonVirtualizationClustersList(),
			"netbox_json_virtualization_interfaces_list": dataNetboxJsonVirtualizationInterfacesList(),
			"netbox_json_virtualization_virtual-machines_list": dataNetboxJsonVirtualizationVirtualMachinesList(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"netbox_ipam_prefix":              resourceNetboxIpamPrefix(),
			"netbox_ipam_ip_addresses":        resourceNetboxIpamIPAddresses(),
			"netbox_ipam_vlan":                resourceNetboxIpamVlan(),
			"netbox_ipam_vlan_group":          resourceNetboxIpamVlanGroup(),
			"netbox_tenancy_tenant":           resourceNetboxTenancyTenant(),
			"netbox_tenancy_tenant_group":     resourceNetboxTenancyTenantGroup(),
			"netbox_virtualization_interface": resourceNetboxVirtualizationInterface(),
			"netbox_virtualization_vm":        resourceNetboxVirtualizationVM(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)
	scheme := d.Get("scheme").(string)

	defaultScheme := []string{scheme}

	t := runtimeclient.New(url, client.DefaultBasePath, defaultScheme)
	t.DefaultAuthentication = runtimeclient.APIKeyAuth(authHeaderName, "header",
		fmt.Sprintf(authHeaderFormat, token))
	return client.New(t, strfmt.Default), nil
}
