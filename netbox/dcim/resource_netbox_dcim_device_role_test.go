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

const resourceNameNetboxDcimDeviceRole = "netbox_dcim_device_role.test"

func TestAccNetboxDcimDeviceRoleMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimDeviceRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimDeviceRoleFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimDeviceRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimDeviceRoleMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
			{
				Config: testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimDeviceRole),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimDeviceRoleConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "dcimdevicerole-{{ .namesuffix }}"
		slug = "dcimdevicerole-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_device_role" "test" {
		name				= "dcimdevicerole-{{ .namesuffix }}"
		slug				= "dcimdevicerole-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test device role"
		color = "00ff00"
		vm_role = false

		{{ if eq .extraresources "true" }}
		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
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
