package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldURL = "netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldURLMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldURL,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldURLFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldURL,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldURLMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldURL),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldURLConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "test_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type          = "url"
		{{ if eq .resourcefull "true" }}
		description   = "Test custom field"
		label         = "Test Label for CF"
		group_name    = "testgroup"
		ui_visibility = "hidden"
		weight        = 50
		#required      = true
		filter_logic  = "disabled"
		default       = jsonencode("https://netbox.dev/")
		{{ end }}
	}

	resource "netbox_dcim_site" "test_assign" {
		name = "test-a-{{ .namesuffix }}"
		slug = "test-a-{{ .namesuffix }}"

		custom_field {
			name = netbox_extras_custom_field.test.name
			type = netbox_extras_custom_field.test.type
			value = "http://my.example.url.invalid/path"
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
