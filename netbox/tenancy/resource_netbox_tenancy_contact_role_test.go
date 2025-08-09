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

const resourceNameNetboxTenancyContactRole = "" +
	"netbox_tenancy_contact_role.test"

func TestAccNetboxTenancyContactRoleMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, false, false, 0),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContactRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactRoleFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContactRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactRoleMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, false, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactRoleConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactRole),
				),
			},
		},
	})
}

func testAccCheckNetboxTenancyContactRoleConfig(nameSuffix string,
	resourceFull, extraResources bool, asn int64) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "tenancycontactrole-{{ .namesuffix }}"
		slug = "tenancycontactrole-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_tenancy_contact_role" "test" {
		name				= "tenancycontactrole-{{ .namesuffix }}"
		slug				= "tenancycontactrole-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test contact role description"

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
		"asn":            strconv.FormatInt(asn, util.Const10),
	}
	return util.RenderTemplate(template, data)
}
