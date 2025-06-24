package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldObject = "" +
	"netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldObjectMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldObject,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldObjectFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldObject,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldObjectMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldObjectConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldObject),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldObjectConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_platform" "test" {
		name = "extrascfobject-{{ .namesuffix }}"
		slug = "extrascfobject-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_extras_custom_field" "test" {
		name = "extrascfobject_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type					= "object"
		object_type	 = "dcim.platform"
		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		group_name		= "testgroup"
		ui_visibility = "hidden"
		ui_editable	 = "no"
		label				 = "Test Label for CF"
		weight				= 50
				#required			= true
		filter_logic	= "disabled"
				{{ if eq .extraresources "true" }}
				default = jsonencode(
						netbox_dcim_platform.test.id
				)
				{{ end }}
		{{ end }}
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test_assign" {
		name = "extrascfobject-{{ .namesuffix }}"
		slug = "extrascfobject-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = netbox_dcim_platform.test.id
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
