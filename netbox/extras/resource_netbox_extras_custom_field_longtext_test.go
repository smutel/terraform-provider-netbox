package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldLongtext = "netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldLongtextMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldLongtext,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldLongtextFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldLongtext,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldLongtextMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldLongtext),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldLongtextConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "test_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type             = "longtext"
		{{ if eq .resourcefull "true" }}
		description      = "Test custom field"
		group_name       = "testgroup"
		ui_visibility    = "hidden"
		label            = "Test Label for CF"
		weight           = 50
		#required         = true
		filter_logic     = "disabled"
		default          = jsonencode("Default text")
		validation_regex = "^.*$"
		{{ end }}
	}

	resource "netbox_dcim_site" "test_assign" {
		name = "test-a-{{ .namesuffix }}"
		slug = "test-a-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = "My text"
		}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"resourcefull":   strconv.FormatBool(resourceFull),
		"extraresources": strconv.FormatBool(extraResources),
	}
	return util.RenderTemplate(template, data)
}
