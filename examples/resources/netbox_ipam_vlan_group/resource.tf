resource "netbox_ipam_vlan_group" "vlan_group_test" {
  name = "TestVlanGroup"
  slug = "TestVlanGroup"

  tag {
    name = "tag1"
    slug = "tag1"
  }
}
