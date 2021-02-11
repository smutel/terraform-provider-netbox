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

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
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
  name    = "Test_VlanGroup"
  slug    = "Test_VlanGroup"
  site_id = data.netbox_dcim_site.site_test.id
}

data "netbox_ipam_role" "vlan_role_production" {
  slug = "production"
}

data "netbox_ipam_role" "vlan_role_backup" {
  slug = "backup"
}

resource "netbox_ipam_vlan" "vlan_test" {
  vlan_id       = 100
  name          = "Test_Vlan"
  site_id       = netbox_ipam_vlan_group.vlan_group_test.site_id
  description   = "VLAN created by terraform"
  vlan_group_id = netbox_ipam_vlan_group.vlan_group_test.id
  tenant_id     = netbox_tenancy_tenant.tenant_test.id
  role_id       = data.netbox_ipam_role.vlan_role_production.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
  }
}

resource "netbox_ipam_prefix" "prefix_test" {
  prefix      = "192.168.56.0/24"
  vlan_id     = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id     = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id     = data.netbox_ipam_role.vlan_role_production.id

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
  address     = "192.168.56.1/24"
  status      = "active"
  tenant_id   = netbox_tenancy_tenant.tenant_test.id
  object_id   = netbox_virtualization_interface.interface_test.id
  object_type = netbox_virtualization_interface.interface_test.type
  primary_ip4 = true

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }
}

data "netbox_virtualization_cluster" "cluster_test" {
  name = "test"
}

data "netbox_dcim_platform" "platform_test" {
  slug = "Debian_10"
}

resource "netbox_virtualization_vm" "vm_test" {
  cluster_id  = data.netbox_virtualization_cluster.cluster_test.id
  name        = "test"
  disk        = 10
  memory      = 10
  platform_id = data.netbox_dcim_platform.platform_test.id
  tenant_id   = netbox_tenancy_tenant.tenant_test.id
  role_id     = 1

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
  }
}

resource "netbox_ipam_service" "service_test" {
  name              = "SMTP"
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  ip_addresses_id   = [netbox_ipam_ip_addresses.ip_test.id]
  port              = "22"
  protocol          = "tcp"
  description       = "Service created by terraform"

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
  }
}

resource "netbox_virtualization_interface" "interface_test" {
  name              = "default"
  virtualmachine_id = netbox_virtualization_vm.vm_test.id
  mac_address       = "AA:AA:AA:AA:AA:AA"
  description       = "Interface de test"
}

resource "netbox_ipam_aggregate" "aggregate_test" {
  prefix      = "192.167.0.0/24"
  rir_id      = 1
  date_added  = "2020-12-21"
  description = "Aggregate created by terraform"

  tag {
    name = "tag1"
    slug = "tag1"
  }

  tag {
    name = "tag2"
    slug = "tag2"
  }

  custom_fields = {
    cf_boolean   = "true"
    cf_date      = "2020-12-25"
    cf_integer   = "10"
    cf_selection = "1"
    cf_text      = "Some text"
    cf_url       = "https://github.com"
  }
}
