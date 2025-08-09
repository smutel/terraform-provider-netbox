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

const resourceNameNetboxDcimSiteGroup = "netbox_dcim_site_group.test"

func TestAccNetboxDcimSiteGroupMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimSiteGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimSiteGroupFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimSiteGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimSiteGroupMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteGroupConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSiteGroup),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimSiteGroupConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "dcimsitegroup-{{ .namesuffix }}"
		slug = "dcimsitegroup-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_site_group" "test" {
		name				= "dcimsitegroup-{{ .namesuffix }}"
		slug				= "dcimsitegroup-{{ .namesuffix }}"
		{{ if eq .resourceFull "true" }}
		description = "Test description"
		{{ end }}

		{{ if eq .extraresources "true" }}
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
