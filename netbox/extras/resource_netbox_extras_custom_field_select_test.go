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

const resourceNameNetboxExtrasCustomFieldSelect = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldSelectMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldSelect,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldSelectFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldSelect,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldSelectMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldSelectConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_extras_custom_field_choice_set" "test" {
		name = "extrascfselect-{{ .namesuffix }}"
		extra_choices {
			value = "extra_choice_1"
			label = "Extra choice 1"
		}

		extra_choices {
			value = "extra_choice_2"
			label = "Extra choice 2"
		}
	}

	resource "netbox_extras_custom_field" "test" {
		name = "extrascfselect_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "select"
		choice_set_name = netbox_extras_custom_field_choice_set.test.name

		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		label				 = "Test Label for CF"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		weight				= 50
		ui_editable	 = "no"
		#required			= true
		filter_logic	= "disabled"
		default			 = jsonencode("extra_choice_1")
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfselect-{{ .namesuffix }}"
		slug = "extrascfselect-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = "extra_choice_2"
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
