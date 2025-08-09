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

const resourceNameNetboxExtrasCustomFieldText = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldTextMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldText,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldTextFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldText,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldTextMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldTextConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldText),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldTextConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "extrascftext_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type						 = "text"
		{{ if eq .resourcefull "true" }}
		description			= "Test custom field"
		label						= "Test Label for CF"
		group_name			 = "testgroup"
		ui_visibility		= "hidden"
		ui_editable			= "no"
		weight					 = 50
		#required				 = true
		filter_logic		 = "disabled"
				default					= jsonencode("Default text")
		validation_regex = "^.*$"
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascftext-{{ .namesuffix }}"
		slug = "extrascftext-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = "My text"
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
