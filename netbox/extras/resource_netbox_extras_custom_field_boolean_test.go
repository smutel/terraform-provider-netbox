package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldBoolean = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldBooleanMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldBoolean,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldBooleanFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldBoolean,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldBooleanMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldBooleanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldBoolean),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldBooleanConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "extrascfboolean_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "boolean"
		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		ui_editable	 = "no"
		label				 = "Test Label for CF"
		weight				= 50
		#required			= true
		filter_logic	= "disabled"
		default			 = true
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfboolean-{{ .namesuffix }}"
		slug = "extrascfboolean-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = false
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
