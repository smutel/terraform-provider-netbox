resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "Test_Tenant"
  slug            = "Test_Tenant"
  description     = "Tenant created by terraform"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }
}

resource "netbox_tenancy_tenant_group" "tenant_group_test" {
  name = "Test_TenantGroup"
  slug = "Test_TenantGroup"
}

data "netbox_dcim_site" "site_test" {
  slug = "pa3"
}

resource "netbox_ipam_vlan_group" "vlan_group_test" {
  name = "Test_VlanGroup"
  slug = "Test_VlanGroup"
  site_id = data.netbox_dcim_site.site_test.id
}

data "netbox_ipam_role" "vlan_role_production" {
  slug = "production"
}

data "netbox_ipam_role" "vlan_role_backup" {
  slug = "backup"
}

resource "netbox_ipam_vlan" "vlan_test" {
  vlan_id = 100
  name = "Test_Vlan"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  description = "VLAN created by terraform"
  vlan_group_id = netbox_ipam_vlan_group.vlan_group_test.id
  tenant_id = netbox_tenancy_tenant.tenant_test.id
  role_id = data.netbox_ipam_role.vlan_role_production.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }
}

resource "netbox_ipam_prefix" "prefix_test" {
  prefix = "192.168.56.0/24"
  vlan_id = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id = data.netbox_ipam_role.vlan_role_production.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }
}

resource "netbox_ipam_ip_addresses" "ip_test" {
  address = "192.168.56.1/24"
  status = "active"
  tenant_id = netbox_tenancy_tenant.tenant_test.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }
}

resource "netbox_virtualization_vm" "vm_test" {
  cluster_id = 1
  name = "test"
  disk = 10
  memory = 10
  platform_id = 1
  tenant_id = netbox_tenancy_tenant.tenant_test.id
  role_id = 1

  tag {
    name = "tag1"
    slug = "tag1"
  }
}

data "netbox_json_dcim_sites_list" "test" {}
output "dcim_sites_list" {
  value = jsondecode(data.netbox_json_dcim_sites_list.test.json)
}
