// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIpamRole = "netbox_ipam_role.test"

func TestAccNetboxIpamRoleMinimal(t *testing.T) {

	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
			{
				ResourceName:      resourceNameIpamRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamRoleFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
			{
				ResourceName:      resourceNameIpamRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamRoleMininmalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
			{
				Config: testAccCheckNetboxIpamRoleConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamRole),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamRoleConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "ipamrole-{{ .namesuffix }}"
		slug = "ipamrole-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_role" "test" {
		name				= "ipamrole-{{ .namesuffix }}"
		slug				= "ipamrole-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test RIR"
		weight = 2000

		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"resourcefull":   strconv.FormatBool(resourceFull),
		"extraresources": strconv.FormatBool(extraResources),
	}
	return util.RenderTemplate(template, data)
}
