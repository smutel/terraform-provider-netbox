package netbox

import (
	"crypto/tls"
	"fmt"
	"net/http"

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
			"basepath": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_BASEPATH", client.DefaultBasePath),
				Description: "URL path to the netbox API.",
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
				Description: "Scheme used to reach netbox application.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_INSECURE", false),
				Description: "Skip TLS certificate validation.",
			},
			"private_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_PRIVATE_KEY_FILE", ""),
				Description: "Private Key used for Secrets",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netbox_dcim_platform":                                dataNetboxDcimPlatform(),
			"netbox_dcim_site":                                    dataNetboxDcimSite(),
			"netbox_ipam_aggregate":                               dataNetboxIpamAggregate(),
			"netbox_ipam_ip_addresses":                            dataNetboxIpamIPAddresses(),
			"netbox_ipam_role":                                    dataNetboxIpamRole(),
			"netbox_ipam_service":                                 dataNetboxIpamService(),
			"netbox_ipam_vlan":                                    dataNetboxIpamVlan(),
			"netbox_ipam_vlan_group":                              dataNetboxIpamVlanGroup(),
			"netbox_json_circuits_circuit_terminations_list":      dataNetboxJSONCircuitsCircuitTerminationsList(),
			"netbox_json_circuits_circuit_types_list":             dataNetboxJSONCircuitsCircuitTypesList(),
			"netbox_json_circuits_circuits_list":                  dataNetboxJSONCircuitsCircuitsList(),
			"netbox_json_circuits_providers_list":                 dataNetboxJSONCircuitsProvidersList(),
			"netbox_json_dcim_cables_list":                        dataNetboxJSONDcimCablesList(),
			"netbox_json_dcim_console_connections_list":           dataNetboxJSONDcimConsoleConnectionsList(),
			"netbox_json_dcim_console_port_templates_list":        dataNetboxJSONDcimConsolePortTemplatesList(),
			"netbox_json_dcim_console_ports_list":                 dataNetboxJSONDcimConsolePortsList(),
			"netbox_json_dcim_console_server_port_templates_list": dataNetboxJSONDcimConsoleServerPortTemplatesList(),
			"netbox_json_dcim_console_server_ports_list":          dataNetboxJSONDcimConsoleServerPortsList(),
			"netbox_json_dcim_device_bay_templates_list":          dataNetboxJSONDcimDeviceBayTemplatesList(),
			"netbox_json_dcim_device_bays_list":                   dataNetboxJSONDcimDeviceBaysList(),
			"netbox_json_dcim_device_roles_list":                  dataNetboxJSONDcimDeviceRolesList(),
			"netbox_json_dcim_device_types_list":                  dataNetboxJSONDcimDeviceTypesList(),
			"netbox_json_dcim_devices_list":                       dataNetboxJSONDcimDevicesList(),
			"netbox_json_dcim_front_port_templates_list":          dataNetboxJSONDcimFrontPortTemplatesList(),
			"netbox_json_dcim_front_ports_list":                   dataNetboxJSONDcimFrontPortsList(),
			"netbox_json_dcim_interface_connections_list":         dataNetboxJSONDcimInterfaceConnectionsList(),
			"netbox_json_dcim_interface_templates_list":           dataNetboxJSONDcimInterfaceTemplatesList(),
			"netbox_json_dcim_interfaces_list":                    dataNetboxJSONDcimInterfacesList(),
			"netbox_json_dcim_inventory_items_list":               dataNetboxJSONDcimInventoryItemsList(),
			"netbox_json_dcim_manufacturers_list":                 dataNetboxJSONDcimManufacturersList(),
			"netbox_json_dcim_platforms_list":                     dataNetboxJSONDcimPlatformsList(),
			"netbox_json_dcim_power_connections_list":             dataNetboxJSONDcimPowerConnectionsList(),
			"netbox_json_dcim_power_feeds_list":                   dataNetboxJSONDcimPowerFeedsList(),
			"netbox_json_dcim_power_outlet_templates_list":        dataNetboxJSONDcimPowerOutletTemplatesList(),
			"netbox_json_dcim_power_outlets_list":                 dataNetboxJSONDcimPowerOutletsList(),
			"netbox_json_dcim_power_panels_list":                  dataNetboxJSONDcimPowerPanelsList(),
			"netbox_json_dcim_power_port_templates_list":          dataNetboxJSONDcimPowerPortTemplatesList(),
			"netbox_json_dcim_power_ports_list":                   dataNetboxJSONDcimPowerPortsList(),
			"netbox_json_dcim_rack_groups_list":                   dataNetboxJSONDcimRackGroupsList(),
			"netbox_json_dcim_rack_reservations_list":             dataNetboxJSONDcimRackReservationsList(),
			"netbox_json_dcim_rack_roles_list":                    dataNetboxJSONDcimRackRolesList(),
			"netbox_json_dcim_racks_list":                         dataNetboxJSONDcimRacksList(),
			"netbox_json_dcim_rear_port_templates_list":           dataNetboxJSONDcimRearPortTemplatesList(),
			"netbox_json_dcim_rear_ports_list":                    dataNetboxJSONDcimRearPortsList(),
			"netbox_json_dcim_regions_list":                       dataNetboxJSONDcimRegionsList(),
			"netbox_json_dcim_sites_list":                         dataNetboxJSONDcimSitesList(),
			"netbox_json_dcim_virtual_chassis_list":               dataNetboxJSONDcimVirtualChassisList(),
			"netbox_json_extras_config_contexts_list":             dataNetboxJSONExtrasConfigContextsList(),
			"netbox_json_extras_export_templates_list":            dataNetboxJSONExtrasExportTemplatesList(),
			"netbox_json_extras_graphs_list":                      dataNetboxJSONExtrasGraphsList(),
			"netbox_json_extras_image_attachments_list":           dataNetboxJSONExtrasImageAttachmentsList(),
			"netbox_json_extras_job_results_list":                 dataNetboxJSONExtrasJobResultsList(),
			"netbox_json_extras_object_changes_list":              dataNetboxJSONExtrasObjectChangesList(),
			"netbox_json_extras_tags_list":                        dataNetboxJSONExtrasTagsList(),
			"netbox_json_ipam_aggregates_list":                    dataNetboxJSONIpamAggregatesList(),
			"netbox_json_ipam_ip_addresses_list":                  dataNetboxJSONIpamIPAddressesList(),
			"netbox_json_ipam_prefixes_list":                      dataNetboxJSONIpamPrefixesList(),
			"netbox_json_ipam_rirs_list":                          dataNetboxJSONIpamRirsList(),
			"netbox_json_ipam_roles_list":                         dataNetboxJSONIpamRolesList(),
			"netbox_json_ipam_services_list":                      dataNetboxJSONIpamServicesList(),
			"netbox_json_ipam_vlan_groups_list":                   dataNetboxJSONIpamVlanGroupsList(),
			"netbox_json_ipam_vlans_list":                         dataNetboxJSONIpamVlansList(),
			"netbox_json_ipam_vrfs_list":                          dataNetboxJSONIpamVrfsList(),
			"netbox_json_secrets_secret_roles_list":               dataNetboxJSONSecretsSecretRolesList(),
			"netbox_json_secrets_secrets_list":                    dataNetboxJSONSecretsSecretsList(),
			"netbox_json_tenancy_tenant_groups_list":              dataNetboxJSONTenancyTenantGroupsList(),
			"netbox_json_tenancy_tenants_list":                    dataNetboxJSONTenancyTenantsList(),
			"netbox_json_users_groups_list":                       dataNetboxJSONUsersGroupsList(),
			"netbox_json_users_permissions_list":                  dataNetboxJSONUsersPermissionsList(),
			"netbox_json_users_users_list":                        dataNetboxJSONUsersUsersList(),
			"netbox_json_virtualization_cluster_groups_list":      dataNetboxJSONVirtualizationClusterGroupsList(),
			"netbox_json_virtualization_cluster_types_list":       dataNetboxJSONVirtualizationClusterTypesList(),
			"netbox_json_virtualization_clusters_list":            dataNetboxJSONVirtualizationClustersList(),
			"netbox_json_virtualization_interfaces_list":          dataNetboxJSONVirtualizationInterfacesList(),
			"netbox_json_virtualization_virtual_machines_list":    dataNetboxJSONVirtualizationVirtualMachinesList(),
			"netbox_tenancy_tenant":                               dataNetboxTenancyTenant(),
			"netbox_tenancy_tenant_group":                         dataNetboxTenancyTenantGroup(),
			"netbox_virtualization_cluster":                       dataNetboxVirtualizationCluster(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"netbox_ipam_aggregate":           resourceNetboxIpamAggregate(),
			"netbox_ipam_ip_addresses":        resourceNetboxIpamIPAddresses(),
			"netbox_ipam_prefix":              resourceNetboxIpamPrefix(),
			"netbox_ipam_service":             resourceNetboxIpamService(),
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
	basepath := d.Get("basepath").(string)
	token := d.Get("token").(string)
	scheme := d.Get("scheme").(string)
	insecure := d.Get("insecure").(bool)
	privateKeyFile := d.Get("private_key_file").(string)

	var options runtimeclient.TLSClientOptions
	options.InsecureSkipVerify = insecure
	tlsConfig, _ := runtimeclient.TLSClientAuth(options)

	headers := make(map[string]string)

	if privateKeyFile != "" {
		privateKey, _ := ReadRSAKey(privateKeyFile)
		sessionKey := GetSessionKey(fmt.Sprintf("%s://%s", scheme, url), token, privateKey)
		headers["X-Session-Key"] = sessionKey
	}
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
