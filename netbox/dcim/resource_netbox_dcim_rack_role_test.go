// Copyright (c)
// SPDX-License-Identifier: MIT

package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxDcimRackRole = "netbox_dcim_rack_role.test"

func TestAccNetboxDcimRackRoleMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRackRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRackRoleFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRackRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRackRoleMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRackRole),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimRackRoleConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "dcimrackrole-{{ .namesuffix }}"
		slug = "dcimrackrole-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_rack_role" "test" {
		name				= "dcimrackrole-{{ .namesuffix }}"
		slug				= "dcimrackrole-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test rack role"
		color = "00ff00"

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
