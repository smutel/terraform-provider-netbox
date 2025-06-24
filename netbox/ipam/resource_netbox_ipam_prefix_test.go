package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIpamPrefix = "netbox_ipam_prefix.test"

func TestAccNetboxIpamPrefixMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
			{
				ResourceName:      resourceNameIpamPrefix,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamPrefixFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
			{
				ResourceName:      resourceNameIpamPrefix,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamPrefixMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
			// This step is necessary. Otherwise deleting site deletes
			// the vlan groups assigned to site.
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
			{
				Config: testAccCheckNetboxIpamPrefixConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamPrefix),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamPrefixConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test" {
		name = "ipamprefix-{{ .namesuffix }}"
		slug = "ipamprefix-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "ipamprefix-{{ .namesuffix }}"
		slug = "ipamprefix-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		 name = "ipamprefix-{{ .namesuffix }}"
		slug = "ipamprefix-{{ .namesuffix }}"
	}

	resource "netbox_ipam_role" "test" {
		 name = "ipamprefix-{{ .namesuffix }}"
		slug = "ipamprefix-{{ .namesuffix }}"
	}

	resource "netbox_ipam_vlan" "test" {
		name				= "ipamprefix-{{ .namesuffix }}"
		vlan_id		 = 10
	}
	{{ end }}

	resource "netbox_ipam_prefix" "test" {
		prefix = "192.168.56.0/24"
		{{ if eq .resourcefull "true" }}
		vlan_id		 = netbox_ipam_vlan.test.id
		description = "Test Vlan group"
		site_id = netbox_dcim_site.test.id
		role_id = netbox_ipam_role.test.id
		tenant_id = netbox_tenancy_tenant.test.id
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
