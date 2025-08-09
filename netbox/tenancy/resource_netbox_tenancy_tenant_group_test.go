// Copyright (c)
// SPDX-License-Identifier: MIT

package tenancy_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxTenancyTenantGroup = "" +
	"netbox_tenancy_tenant_group.test"

func TestAccNetboxTenancyTenantGroupMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyTenantGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyTenantGroupFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyTenantGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyTenantGroupMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyTenantGroup),
				),
			},
		},
	})
}

func testAccCheckNetboxTenancyTenantGroupConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "tenancytenantgroup-{{ .namesuffix }}"
		slug = "tenancytenantgroup-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_tenancy_tenant_group" "test" {
		name				= "tenancytenantgroup-{{ .namesuffix }}"
		slug				= "tenancytenantgroup-{{ .namesuffix }}"
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
