package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldMultiSelect = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldMultiSelectMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldMultiSelect, //nolint:revive
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldMultiSelectFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldMultiSelect, //nolint:revive
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldMultiSelectMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiSelect),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldMultiSelectConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_extras_custom_field_choice_set" "test" {
		name = "extrascfmultiselect-{{ .namesuffix }}"
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
		name = "extrascfmultiselect_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "multiselect"
		choice_set_name = netbox_extras_custom_field_choice_set.test.name

		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		label				 = "Test Label for CF"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		ui_editable	 = "no"
		weight				= 50
		#required			= true
		filter_logic	= "disabled"
		default			 = jsonencode(["extra_choice_1"])
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfmultiselect-{{ .namesuffix }}"
		slug = "extrascfmultiselect-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = jsonencode(
				[
					"extra_choice_1",
				]
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
