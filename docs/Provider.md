# Netbox Provider

Table of Contents
=================

   * [Netbox Provider](#netbox-provider)
      * [Data Sources](#data-sources)
         * [netbox_dcim_site](#netbox_dcim_site)
         * [netbox_ipam_role](#netbox_ipam_role)
         * [netbox_ipam_vlan_group](#netbox_ipam_vlan_group)
         * [netbox_ipam_vlan](#netbox_ipam_vlan)
         * [netbox_tenancy_tenant_group](#netbox_tenancy_tenant_group)
         * [netbox_tenancy_tenant](#netbox_tenancy_tenant)
      * [Resources](#resources)
         * [netbox_ipam_prefix](#netbox_ipam_prefix)
         * [netbox_ipam_vlan](#netbox_ipam_vlan-1)
         * [netbox_ipam_vlan_group](#netbox_ipam_vlan_group-1)
         * [netbox_tenancy_tenant](#netbox_tenancy_tenant-1)
         * [netbox_tenancy_tenant_group](#netbox_tenancy_tenant_group-1)
   * [Table of Contents](#table-of-contents)

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
* ``id`` - The id of this dcim site.

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
* ``id`` - The id of this ipam role.

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
* ``id`` - The id of this ipam vlan groups.

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
* ``id`` - The id of this ipam vlan.

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
* ``id`` - The id of this tenant group.

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
* ``id`` - The id of this tenant.

## Resources

### netbox_ipam_prefix

Manages an ipam prefix resource within Netbox.


__**Example Usage**__

```hcl
resource "netbox_pam_prefix" "prefix_test" {
  prefix = "192.168.56.0/24"
  vlan_id = netbox_ipam_vlan.vlan_test.id
  description = "Prefix created by terraform"
  site_id = netbox_ipam_vlan_group.vlan_group_test.site_id
  role_id = data.netbox_ipam_roles.vlan_role_production.id
  tags = ["tag1"]
  status = "Active"
}
```

__**Argument Reference**__

The following arguments are supported:
* ``prefix`` - (Required) The prefix (IP address/mask) used for this object.
* ``status`` - (Optional) The status among Container, Active, Reserved, Deprecated (Active by default).
* ``vrf_id`` - (Optional) The ID of the vrf attached to this prefix.
* ``role_id`` - (Optional) The ID of the role attached to this prefix.
* ``description`` - (Optional) The description of this prefix.
* ``is_pool`` - (Optional) Define if this prefix is a pool (false by default).
* ``site_id`` - (Optional) ID of the site where this prefix is created
* ``vlan_id`` - (Optional) ID of the vlan where this prefix is attached.
* ``tenant_id`` - (Optional) ID of the tenant where this prefix is attached.
* ``custom_field`` - (Optional) List of custom fields.
** ``name`` - (Required) Name of the custom field.
** ``kind`` - (Required) Type of custom field among string, bool and int. Custom
field of type Choice are not available in this provider.
** ``value`` - (Required) Value for this custom field.
* ``tags`` - (Optional) Array of tags for this prefix.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id of this prefix.

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
* ``vlan_id`` - (Required) The ID of this vlan.
* ``status`` - (Optional) The status among Container, Active, Reserved, Deprecated (Active by default).
* ``site_id`` - (Optional) ID of the site where this vlan is created.
* ``vlan_group_id`` - (Optional) ID of the group where this vlan belongs to.
* ``role_id`` - (Optional) The ID of the role attached to this vlan.
* ``description`` - (Optional) The description of this vlan.
* ``tenant_id`` - (Optional) ID of the tenant where this vlan is attached.
* ``custom_field`` - (Optional) List of custom fields.
** ``name`` - (Required) Name of the custom field.
** ``kind`` - (Required) Type of custom field among string, bool and int. Custom
field of type Choice are not available in this provider.
** ``value`` - (Required) Value for this custom field.
* ``tags`` - (Optional) Array of tags for this vlan.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id of this vlan.

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
* ``name`` - (Required) The name for this vlan group.
* ``slug`` - (Required) The slug for this vlan group.
* ``site_id`` - (Optional) ID of the site where this vlan group is created.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id of this vlan group.

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
* ``name`` - (Required) The name for this tenant.
* ``slug`` - (Required) The slug for this tenant.
* ``description`` - (Optional) The description for this tenant.
* ``comments`` - (Optional) Comments for this tenant.
* ``tenant_group_id`` - (Optional) ID of the group where this tenant is located.
* ``custom_field`` - (Optional) List of custom fields.
** ``name`` - (Required) Name of the custom field.
** ``kind`` - (Required) Type of custom field among string, bool and int. Custom
field of type Choice are not available in this provider.
** ``value`` - (Required) Value for this custom field.
* ``tags`` - (Optional) Array of tags for this tenant.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id of this tenant.

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
* ``name`` - (Required) The name for this tenant group.
* ``slug`` - (Required) The slug for this tenant group.

__**Attributes Reference**__

In addition to the above arguments, the following attributes are exported:
* ``id`` - The id of this tenant group.
