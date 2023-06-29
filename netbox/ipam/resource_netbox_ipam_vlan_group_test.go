package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

const resourceNameVlanGroup = "netbox_ipam_vlan_group.test"

func TestAccNetboxIpamVlanGroupMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
			{
				ResourceName:      resourceNameVlanGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVlanGroupFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
			{
				ResourceName:      resourceNameVlanGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVlanGroupMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
			// This step is necessary. Otherwise deleting site deletes the vlan groups assigned to site.
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
			{
				Config: testAccCheckNetboxIpamVlanGroupConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameVlanGroup),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamVlanGroupConfig(nameSuffix string, resourceFull, extraResources bool) string {
	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_vlan_group" "test" {
		name        = "test-{{ .namesuffix }}"
		slug        = "test-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test Vlan group"
		max_vid     = 2369
		min_vid     = 135

		scope {
		  id = netbox_dcim_site.test.id
		  type = "dcim.site"
		}

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
