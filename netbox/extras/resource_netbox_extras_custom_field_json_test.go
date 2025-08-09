// Copyright (c)
// SPDX-License-Identifier: MIT

package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldJSON = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldJSONMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldJSON,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldJSONFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldJSON,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldJSONMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldJSONConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldJSON),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldJSONConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "extrascfjson_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "json"
		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		ui_editable	 = "no"
		label				 = "Test Label for CF"
		weight				= 50
		#required			= true
		filter_logic	= "disabled"
		default			 = jsonencode({
			bool = false
			number = 1.5
			dict = {
				text = "Some text"
			}
		})
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfjson-a-{{ .namesuffix }}"
		slug = "extrascfjson-a-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = jsonencode(
				{
					mystring = "string"
					mynumber = 6
					mydict = {
						mybool = false
					}
				}
			)
		}
	}
	{{ end }}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"resourcefull":   strconv.FormatBool(resourceFull),
		"extraresources": strconv.FormatBool(extraResources),
	}
	return util.RenderTemplate(template, data)
}
