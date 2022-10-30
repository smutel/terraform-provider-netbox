package netbox

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/dcim"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/extras"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/ipam"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/json"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/tenancy"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/virtualization"
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
				Description: "URL and port to reach netbox application (127.0.0.1:8000 by default).",
			},
			"basepath": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_BASEPATH", client.DefaultBasePath),
				Description: "URL base path to the netbox API (/api by default).",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_TOKEN", ""),
				Description: "Token used for API operations (empty by default).",
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_SCHEME", "https"),
				Description: "Scheme used to reach netbox application (https by default).",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_INSECURE", false),
				Description: "Skip TLS certificate validation (false by default).",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netbox_json_circuits_circuits_list":                  json.DataNetboxJSONCircuitsCircuitsList(),
			"netbox_json_circuits_circuit_terminations_list":      json.DataNetboxJSONCircuitsCircuitTerminationsList(),
			"netbox_json_circuits_circuit_types_list":             json.DataNetboxJSONCircuitsCircuitTypesList(),
			"netbox_json_circuits_provider_networks_list":         json.DataNetboxJSONCircuitsProviderNetworksList(),
			"netbox_json_circuits_providers_list":                 json.DataNetboxJSONCircuitsProvidersList(),
			"netbox_json_dcim_cables_list":                        json.DataNetboxJSONDcimCablesList(),
			"netbox_json_dcim_console_ports_list":                 json.DataNetboxJSONDcimConsolePortsList(),
			"netbox_json_dcim_console_port_templates_list":        json.DataNetboxJSONDcimConsolePortTemplatesList(),
			"netbox_json_dcim_console_server_ports_list":          json.DataNetboxJSONDcimConsoleServerPortsList(),
			"netbox_json_dcim_console_server_port_templates_list": json.DataNetboxJSONDcimConsoleServerPortTemplatesList(),
			"netbox_json_dcim_device_bays_list":                   json.DataNetboxJSONDcimDeviceBaysList(),
			"netbox_json_dcim_device_bay_templates_list":          json.DataNetboxJSONDcimDeviceBayTemplatesList(),
			"netbox_json_dcim_device_roles_list":                  json.DataNetboxJSONDcimDeviceRolesList(),
			"netbox_json_dcim_devices_list":                       json.DataNetboxJSONDcimDevicesList(),
			"netbox_json_dcim_device_types_list":                  json.DataNetboxJSONDcimDeviceTypesList(),
			"netbox_json_dcim_front_ports_list":                   json.DataNetboxJSONDcimFrontPortsList(),
			"netbox_json_dcim_front_port_templates_list":          json.DataNetboxJSONDcimFrontPortTemplatesList(),
			"netbox_json_dcim_interfaces_list":                    json.DataNetboxJSONDcimInterfacesList(),
			"netbox_json_dcim_interface_templates_list":           json.DataNetboxJSONDcimInterfaceTemplatesList(),
			"netbox_json_dcim_inventory_item_roles_list":          json.DataNetboxJSONDcimInventoryItemRolesList(),
			"netbox_json_dcim_inventory_items_list":               json.DataNetboxJSONDcimInventoryItemsList(),
			"netbox_json_dcim_inventory_item_templates_list":      json.DataNetboxJSONDcimInventoryItemTemplatesList(),
			"netbox_json_dcim_locations_list":                     json.DataNetboxJSONDcimLocationsList(),
			"netbox_json_dcim_manufacturers_list":                 json.DataNetboxJSONDcimManufacturersList(),
			"netbox_json_dcim_module_bays_list":                   json.DataNetboxJSONDcimModuleBaysList(),
			"netbox_json_dcim_module_bay_templates_list":          json.DataNetboxJSONDcimModuleBayTemplatesList(),
			"netbox_json_dcim_modules_list":                       json.DataNetboxJSONDcimModulesList(),
			"netbox_json_dcim_module_types_list":                  json.DataNetboxJSONDcimModuleTypesList(),
			"netbox_json_dcim_platforms_list":                     json.DataNetboxJSONDcimPlatformsList(),
			"netbox_json_dcim_power_feeds_list":                   json.DataNetboxJSONDcimPowerFeedsList(),
			"netbox_json_dcim_power_outlets_list":                 json.DataNetboxJSONDcimPowerOutletsList(),
			"netbox_json_dcim_power_outlet_templates_list":        json.DataNetboxJSONDcimPowerOutletTemplatesList(),
			"netbox_json_dcim_power_panels_list":                  json.DataNetboxJSONDcimPowerPanelsList(),
			"netbox_json_dcim_power_ports_list":                   json.DataNetboxJSONDcimPowerPortsList(),
			"netbox_json_dcim_power_port_templates_list":          json.DataNetboxJSONDcimPowerPortTemplatesList(),
			"netbox_json_dcim_rack_reservations_list":             json.DataNetboxJSONDcimRackReservationsList(),
			"netbox_json_dcim_rack_roles_list":                    json.DataNetboxJSONDcimRackRolesList(),
			"netbox_json_dcim_racks_list":                         json.DataNetboxJSONDcimRacksList(),
			"netbox_json_dcim_rear_ports_list":                    json.DataNetboxJSONDcimRearPortsList(),
			"netbox_json_dcim_rear_port_templates_list":           json.DataNetboxJSONDcimRearPortTemplatesList(),
			"netbox_json_dcim_regions_list":                       json.DataNetboxJSONDcimRegionsList(),
			"netbox_json_dcim_site_groups_list":                   json.DataNetboxJSONDcimSiteGroupsList(),
			"netbox_json_dcim_sites_list":                         json.DataNetboxJSONDcimSitesList(),
			"netbox_json_dcim_virtual_chassis_list":               json.DataNetboxJSONDcimVirtualChassisList(),
			"netbox_json_extras_config_contexts_list":             json.DataNetboxJSONExtrasConfigContextsList(),
			"netbox_json_extras_content_types_list":               json.DataNetboxJSONExtrasContentTypesList(),
			"netbox_json_extras_custom_fields_list":               json.DataNetboxJSONExtrasCustomFieldsList(),
			"netbox_json_extras_custom_links_list":                json.DataNetboxJSONExtrasCustomLinksList(),
			"netbox_json_extras_export_templates_list":            json.DataNetboxJSONExtrasExportTemplatesList(),
			"netbox_json_extras_image_attachments_list":           json.DataNetboxJSONExtrasImageAttachmentsList(),
			"netbox_json_extras_job_results_list":                 json.DataNetboxJSONExtrasJobResultsList(),
			"netbox_json_extras_journal_entries_list":             json.DataNetboxJSONExtrasJournalEntriesList(),
			"netbox_json_extras_object_changes_list":              json.DataNetboxJSONExtrasObjectChangesList(),
			"netbox_json_extras_tags_list":                        json.DataNetboxJSONExtrasTagsList(),
			"netbox_json_extras_webhooks_list":                    json.DataNetboxJSONExtrasWebhooksList(),
			"netbox_json_ipam_aggregates_list":                    json.DataNetboxJSONIpamAggregatesList(),
			"netbox_json_ipam_asns_list":                          json.DataNetboxJSONIpamAsnsList(),
			"netbox_json_ipam_fhrp_group_assignments_list":        json.DataNetboxJSONIpamFhrpGroupAssignmentsList(),
			"netbox_json_ipam_fhrp_groups_list":                   json.DataNetboxJSONIpamFhrpGroupsList(),
			"netbox_json_ipam_ip_addresses_list":                  json.DataNetboxJSONIpamIPAddressesList(),
			"netbox_json_ipam_ip_ranges_list":                     json.DataNetboxJSONIpamIPRangesList(),
			"netbox_json_ipam_prefixes_list":                      json.DataNetboxJSONIpamPrefixesList(),
			"netbox_json_ipam_rirs_list":                          json.DataNetboxJSONIpamRirsList(),
			"netbox_json_ipam_roles_list":                         json.DataNetboxJSONIpamRolesList(),
			"netbox_json_ipam_route_targets_list":                 json.DataNetboxJSONIpamRouteTargetsList(),
			"netbox_json_ipam_services_list":                      json.DataNetboxJSONIpamServicesList(),
			"netbox_json_ipam_service_templates_list":             json.DataNetboxJSONIpamServiceTemplatesList(),
			"netbox_json_ipam_vlan_groups_list":                   json.DataNetboxJSONIpamVlanGroupsList(),
			"netbox_json_ipam_vlans_list":                         json.DataNetboxJSONIpamVlansList(),
			"netbox_json_ipam_vrfs_list":                          json.DataNetboxJSONIpamVrfsList(),
			"netbox_json_tenancy_contact_assignments_list":        json.DataNetboxJSONTenancyContactAssignmentsList(),
			"netbox_json_tenancy_contact_groups_list":             json.DataNetboxJSONTenancyContactGroupsList(),
			"netbox_json_tenancy_contact_roles_list":              json.DataNetboxJSONTenancyContactRolesList(),
			"netbox_json_tenancy_contacts_list":                   json.DataNetboxJSONTenancyContactsList(),
			"netbox_json_tenancy_tenant_groups_list":              json.DataNetboxJSONTenancyTenantGroupsList(),
			"netbox_json_tenancy_tenants_list":                    json.DataNetboxJSONTenancyTenantsList(),
			"netbox_json_users_groups_list":                       json.DataNetboxJSONUsersGroupsList(),
			"netbox_json_users_permissions_list":                  json.DataNetboxJSONUsersPermissionsList(),
			"netbox_json_users_tokens_list":                       json.DataNetboxJSONUsersTokensList(),
			"netbox_json_users_users_list":                        json.DataNetboxJSONUsersUsersList(),
			"netbox_json_virtualization_cluster_groups_list":      json.DataNetboxJSONVirtualizationClusterGroupsList(),
			"netbox_json_virtualization_clusters_list":            json.DataNetboxJSONVirtualizationClustersList(),
			"netbox_json_virtualization_cluster_types_list":       json.DataNetboxJSONVirtualizationClusterTypesList(),
			"netbox_json_virtualization_interfaces_list":          json.DataNetboxJSONVirtualizationInterfacesList(),
			"netbox_json_virtualization_virtual_machines_list":    json.DataNetboxJSONVirtualizationVirtualMachinesList(),
			"netbox_json_wireless_wireless_lan_groups_list":       json.DataNetboxJSONWirelessWirelessLanGroupsList(),
			"netbox_json_wireless_wireless_lans_list":             json.DataNetboxJSONWirelessWirelessLansList(),
			"netbox_json_wireless_wireless_links_list":            json.DataNetboxJSONWirelessWirelessLinksList(),
			"netbox_dcim_platform":                                dcim.DataNetboxDcimPlatform(),
			"netbox_dcim_site":                                    dcim.DataNetboxDcimSite(),
			"netbox_ipam_aggregate":                               ipam.DataNetboxIpamAggregate(),
			"netbox_ipam_ip_addresses":                            ipam.DataNetboxIpamIPAddresses(),
			"netbox_ipam_role":                                    ipam.DataNetboxIpamRole(),
			"netbox_ipam_service":                                 ipam.DataNetboxIpamService(),
			"netbox_ipam_vlan":                                    ipam.DataNetboxIpamVlan(),
			"netbox_ipam_vlan_group":                              ipam.DataNetboxIpamVlanGroup(),
			"netbox_tenancy_contact":                              tenancy.DataNetboxTenancyContact(),
			"netbox_tenancy_contact_group":                        tenancy.DataNetboxTenancyContactGroup(),
			"netbox_tenancy_contact_role":                         tenancy.DataNetboxTenancyContactRole(),
			"netbox_tenancy_tenant":                               tenancy.DataNetboxTenancyTenant(),
			"netbox_tenancy_tenant_group":                         tenancy.DataNetboxTenancyTenantGroup(),
			"netbox_virtualization_cluster":                       virtualization.DataNetboxVirtualizationCluster(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"netbox_dcim_manufacturer":            dcim.ResourceNetboxDcimManufacturer(),
			"netbox_dcim_device_role":             dcim.ResourceNetboxDcimDeviceRole(),
			"netbox_dcim_platform":                dcim.ResourceNetboxDcimPlatform(),
			"netbox_dcim_site":                    dcim.ResourceNetboxDcimSite(),
			"netbox_extras_tag":                   extras.ResourceNetboxExtrasTag(),
			"netbox_ipam_aggregate":               ipam.ResourceNetboxIpamAggregate(),
			"netbox_ipam_asn":                     ipam.ResourceNetboxIpamASN(),
			"netbox_ipam_ip_addresses":            ipam.ResourceNetboxIpamIPAddresses(),
			"netbox_ipam_ip_range":                ipam.ResourceNetboxIpamIPRange(),
			"netbox_ipam_prefix":                  ipam.ResourceNetboxIpamPrefix(),
			"netbox_ipam_rir":                     ipam.ResourceNetboxIpamRIR(),
			"netbox_ipam_service":                 ipam.ResourceNetboxIpamService(),
			"netbox_ipam_vlan":                    ipam.ResourceNetboxIpamVlan(),
			"netbox_ipam_vlan_group":              ipam.ResourceNetboxIpamVlanGroup(),
			"netbox_tenancy_contact":              tenancy.ResourceNetboxTenancyContact(),
			"netbox_tenancy_contact_assignment":   tenancy.ResourceNetboxTenancyContactAssignment(),
			"netbox_tenancy_contact_group":        tenancy.ResourceNetboxTenancyContactGroup(),
			"netbox_tenancy_contact_role":         tenancy.ResourceNetboxTenancyContactRole(),
			"netbox_tenancy_tenant":               tenancy.ResourceNetboxTenancyTenant(),
			"netbox_tenancy_tenant_group":         tenancy.ResourceNetboxTenancyTenantGroup(),
			"netbox_virtualization_interface":     virtualization.ResourceNetboxVirtualizationInterface(),
			"netbox_virtualization_vm":            virtualization.ResourceNetboxVirtualizationVM(),
			"netbox_virtualization_vm_primary_ip": virtualization.ResourceNetboxVirtualizationVMPrimaryIP(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	basepath := d.Get("basepath").(string)
	token := d.Get("token").(string)
	scheme := d.Get("scheme").(string)
	insecure := d.Get("insecure").(bool)

	var options runtimeclient.TLSClientOptions
	options.InsecureSkipVerify = insecure
	tlsConfig, _ := runtimeclient.TLSClientAuth(options)

	headers := make(map[string]string)

	// Create a custom client
	// Override the default transport with a RoundTripper to inject dynamic headers
	// Add TLSOptions
	cli := &http.Client{
		Transport: &transport{
			headers:         headers,
			TLSClientConfig: tlsConfig,
		},
	}

	defaultScheme := []string{scheme}

	t := runtimeclient.NewWithClient(url, basepath, defaultScheme, cli)
	t.DefaultAuthentication = runtimeclient.APIKeyAuth(authHeaderName, "header",
		fmt.Sprintf(authHeaderFormat, token))

	return client.New(t, strfmt.Default), nil
}

type transport struct {
	headers         map[string]string
	base            http.RoundTripper
	TLSClientConfig *tls.Config
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add headers to request
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}
	base := t.base
	if base == nil {
		// init an http.Transport with TLSOptions
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = t.TLSClientConfig
		base = customTransport
	}
	return base.RoundTrip(req)
}
