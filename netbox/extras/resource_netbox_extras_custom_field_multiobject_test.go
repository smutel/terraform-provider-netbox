package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldMultiObject = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldMultiObjectMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldMultiObject, //nolint:revive
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldMultiObjectFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldMultiObject, //nolint:revive
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldMultiObjectMinimalFullMinimal(
	t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldMultiObject),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldMultiObjectConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_platform" "test" {
		name = "extrascfmultiobject-{{ .namesuffix }}"
		slug = "extrascfmultiobject-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_extras_custom_field" "test" {
		name = "extrascfmultiobject_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "multiobject"
		object_type	 = "dcim.platform"
		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		label				 = "Test Label for CF"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		ui_editable	 = "no"
		weight				= 50
		#required			= true
		filter_logic	= "disabled"
		{{ if eq .extraresources "true" }}
		default			 = jsonencode([
			netbox_dcim_platform.test.id
		])
		{{ end }}
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfmultiobject-{{ .namesuffix }}"
		slug = "extrascfmultiobject-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = jsonencode(
				[
					netbox_dcim_platform.test.id,
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
