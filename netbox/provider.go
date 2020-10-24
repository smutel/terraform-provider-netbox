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
