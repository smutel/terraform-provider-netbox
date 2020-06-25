# Netbox Provider

Table of Contents
=================

   * [Netbox Provider](#netbox-provider)
      * [Data Sources](#data-sources)
         * [netbox_dcim_site](#netbox_dcim_site)
         * [netbox_ipam_ip_addresses](#netbox_ipam_ip_addresses)
         * [netbox_ipam_role](#netbox_ipam_role)
         * [netbox_ipam_vlan](#netbox_ipam_vlan)
         * [netbox_ipam_vlan_group](#netbox_ipam_vlan_group)
         * [netbox_tenancy_tenant](#netbox_tenancy_tenant)
         * [netbox_tenancy_tenant_group](#netbox_tenancy_tenant_group)
      * [Resources](#resources)
         * [netbox_ipam_ip_addresses](#netbox_ipam_ip_addresses-1)
         * [netbox_ipam_prefix](#netbox_ipam_prefix-1)
         * [netbox_ipam_vlan](#netbox_ipam_vlan-1)
         * [netbox_ipam_vlan_group](#netbox_ipam_vlan_group-1)
         * [netbox_tenancy_tenant](#netbox_tenancy_tenant-1)
         * [netbox_tenancy_tenant_group](#netbox_tenancy_tenant_group-1)

## Data Sources

### netbox_dcim_site

Get info about dcim site in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_dcim_site" "site_test" {
  slug = "TestSite"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``slug`` - (Required) The slug of the dcim site.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_ip_addresses

Get info about ipam IP addresses in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_ipam_ip_addresses" "ipaddress_test" {
  address = "192.168.56.1/24"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``address`` - (Required) The address (with mask) of the ipam IP address.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_role

Get info about ipam role in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_ipam_role" "role_test" {
  slug = "TestRole"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``slug`` - (Required) The slug of the ipam role.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_vlan_group

Get info about ipam vlan group in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_ipam_vlan_group" "vlan_group_test" {
  slug = "TestVlanGroup"
  site_id = 15
}
```

__**Argument Reference**__

The following arguments are supported:
* ``slug`` - (Required) The slug of the ipam vlan group.
* ``site_id`` - (Optional) The site_id of the ipam vlan groups.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_vlan

Get info about ipam vlan in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_ipam_vlan" "vlan_test" {
  vlan_id = 15
  vlan_group_id = 16
}
```

__**Argument Reference**__

The following arguments are supported:
* ``vlan_id`` - (Required) The ID of the ipam vlan.
* ``vlan_group_id`` - (Optional) The id of the vlan group linked to this vlan.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_tenancy_tenant_group

Get info about tenancy tenant groups in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_tenancy_tenant_groups" "tenant_group_test" {
  slug = "TestTenantGroup"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenancy tenant groups.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_tenancy_tenant

Get info about tenancy tenant in the netbox provider.

__**Example Usage**__

```hcl
data "netbox_tenancy_tenant" "tenant_test" {
  slug = "TestTenant"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``slug`` - (Required) The slug of the tenant.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

## Resources

### netbox_ipam_ip_addresses

Manages an ipam ip addresses resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_ipam_ip_addresses" "ip_test" {
  address = "192.168.56.0/24"
  description = "IP created by terraform"
  tags = ["tag1"]
  status = "active"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``address`` - (Required) The IP address (with mask) used for this object.
* ``description`` - (Optional) The description of this object.
* ``dns_name`` - (Optional) The DNS name of this object.
* ``interface_id`` - (Optional) The ID of the interface where this object is attached to.
* ``nat_inside_id`` - (Optional) The ID of the NAT inside of this object.
* ``nat_outside_id`` - (Optional) The ID of the NAT outside of this object.
* ``role`` - (Optional) The role among loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp of this object.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tags`` - (Optional) Array of tags for this object.
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_prefix

Manages an ipam prefix resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_ipam_prefix" "prefix_test" {
  prefix = "192.168.56.0/24"
  vlan_id = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
  tags = ["tag1"]
  status = "active"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``description`` - (Optional) The description of this object.
* ``is_pool`` - (Optional) Define if this object is a pool (false by default).
* ``prefix`` - (Required) The prefix (IP address/mask) used for this object.
* ``role_id`` - (Optional) The ID of the role attached to this object.
* ``site_id`` - (Optional) ID of the site where this object is created
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tags`` - (Optional) Array of tags for this object.
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vlan_id`` - (Optional) ID of the vlan where this object is attached.
* ``vrf_id`` - (Optional) The ID of the vrf attached to this object.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_vlan

Manages an ipam vlan resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_ipam_vlan" "vlan_test" {
  vlan_id = 100
  name = "TestVlan"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  description = "VLAN created by terraform"
  vlan_group_id = netbox_ipam_vlan_group.vlan_group_test.id
  tenant_id = netbox_tenancy_tenant.tenant_test.id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
  tags = ["tag1"]
}
```

__**Argument Reference**__

The following arguments are supported:
* ``description`` - (Optional) The description of this object.
* ``vlan_group_id`` - (Optional) ID of the group where this object belongs to.
* ``name`` - (Required) The name for this object.
* ``role_id`` - (Optional) The ID of the role attached to this object.
* ``site_id`` - (Optional) ID of the site where this object is created.
* ``status`` - (Optional) The status among container, active, reserved, deprecated (active by default).
* ``tags`` - (Optional) Array of tags for this vlan.
* ``tenant_id`` - (Optional) ID of the tenant where this object is attached.
* ``vlan_id`` - (Required) The ID of this vlan (vlan tag).

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_ipam_vlan_group

Manages an ipam vlan group resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_ipam_vlan_group" "vlan_group_test" {
  name = "TestVlanGroup"
  slug = "TestVlanGroup"
  site_id = data.netbox_dcim_sites.site_test.id
}
```

__**Argument Reference**__

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``site_id`` - (Optional) ID of the site where this object is created.
* ``slug`` - (Required) The slug for this object.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_tenancy_tenant

Manages an tenancy tenant resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "TestTenant"
  slug            = "TestTenant"
  description     = "Tenant created by terraform"
  comments        = "Some test comments"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id
  tags            = ["tag1"]
}
```

__**Argument Reference**__

The following arguments are supported:
* ``comments`` - (Optional) Comments for this object.
* ``description`` - (Optional) The description for this object.
* ``tenant_group_id`` - (Optional) ID of the group where this object is located.
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.
* ``tags`` - (Optional) Array of tags for this tenant.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.

### netbox_tenancy_tenant_group

Manages an tenancy tenant group resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_tenancy_tenant_group" "tenant_group_test" {
  name = "TestTenantGroup"
  slug = "TestTenantGroup"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``name`` - (Required) The name for this object.
* ``slug`` - (Required) The slug for this object.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id (ref in Netbox) of this object.
