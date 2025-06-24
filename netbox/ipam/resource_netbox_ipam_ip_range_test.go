package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIpamIPRange = "netbox_ipam_ip_range.test"

func TestAccNetboxIpamIPRangeMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
			{
				ResourceName:      resourceNameIpamIPRange,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamIPRangeFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
			{
				ResourceName:      resourceNameIpamIPRange,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamIPRangeMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
			// This step is necessary. Otherwise deleting site deletes
			// the vlan groups assigned to site.
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
			{
				Config: testAccCheckNetboxIpamIPRangeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamIPRange),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamIPRangeConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "ipamiprange-{{ .namesuffix }}"
		slug = "ipamiprange-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		 name = "ipamiprange-{{ .namesuffix }}"
		slug = "ipamiprange-{{ .namesuffix }}"
	}

	resource "netbox_ipam_role" "test" {
		 name = "ipamiprange-{{ .namesuffix }}"
		slug = "ipamiprange-{{ .namesuffix }}"
	}

	resource "netbox_ipam_vrf" "test" {
		name = "ipamiprange-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_ip_range" "test" {
		start_address = "192.168.56.1/24"
		end_address = "192.168.56.100/24"
		{{ if eq .resourcefull "true" }}
		description = "Test IP range"
		role_id = netbox_ipam_role.test.id
		tenant_id = netbox_tenancy_tenant.test.id
		vrf_id = netbox_ipam_vrf.test.id
		status = "reserved"

		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"extraresources": strconv.FormatBool(extraResources),
		"resourcefull":   strconv.FormatBool(resourceFull),
	}
	return util.RenderTemplate(template, data)
}
