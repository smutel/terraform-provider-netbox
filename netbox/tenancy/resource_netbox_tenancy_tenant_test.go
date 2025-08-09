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

const resourceNameNetboxTenancyTenant = "netbox_tenancy_tenant.test"

func TestAccNetboxTenancyTenantMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, false, false, 0),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyTenant,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyTenantFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyTenant,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyTenantMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, false, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
			{
				Config: testAccCheckNetboxTenancyTenantConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxTenancyTenant),
				),
			},
		},
	})
}

func testAccCheckNetboxTenancyTenantConfig(nameSuffix string,
	resourceFull, extraResources bool, asn int64) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_tenancy_tenant_group" "tenant_group_test" {
		name = "tenancytenant-{{ .namesuffix }}"
		slug = "tenancytenant-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "tenancytenant-{{ .namesuffix }}"
		slug = "tenancytenant-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_tenancy_tenant" "test" {
		name				= "tenancytenant-{{ .namesuffix }}"
		slug				= "tenancytenant-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test tenant description"
		comments		= "Test tenant comments"

		{{ if eq .extraresources "true" }}
		tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id

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
		"asn":            strconv.FormatInt(asn, util.Const10),
	}
	return util.RenderTemplate(template, data)
}
