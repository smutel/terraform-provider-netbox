resource "netbox_tenancy_tenant_group" "tenant_group_test" {
  name = "TestTenantGroup"
  slug = "TestTenantGroup"

  tag {
    name = "tag1"
    slug = "tag1"
  }
}
