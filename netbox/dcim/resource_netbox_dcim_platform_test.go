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

const resourceNameNetboxDcimPlatform = "netbox_dcim_platform.test"

func TestAccNetboxDcimPlatformMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimPlatform,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimPlatformFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimPlatform,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimPlatformMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
			{
				Config: testAccCheckNetboxDcimPlatformConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimPlatform),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimPlatformConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_manufacturer" "test" {
		name = "dcimplatform-{{ .namesuffix }}"
		slug = "dcimplatform-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "dcimplatform-{{ .namesuffix }}"
		slug = "dcimplatform-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_platform" "test" {
		name						= "dcimplatform-{{ .namesuffix }}"
		slug						= "dcimplatform-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description		 = "Test device role"
		manufacturer_id = netbox_dcim_manufacturer.test.id

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
